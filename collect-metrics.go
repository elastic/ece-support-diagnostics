package ecediag

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	elasticsearch "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

// Need to refactor into modules. Discovery should set username/password into a common.Config?

func ScrollRunner(ece *ECEendpoint, tar *Tarball) {
	tp := CustomRoundTripper{
		RoundTripper: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				// MinVersion:         tls.VersionTLS12,
			},
			// Connection timeout = 5s
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			// TLS Handshake Timeout = 5s
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	metricCluster, err := ece.discoverMetricsCluster(&tp)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(ece.eceAPI)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "/api/v1/clusters/elasticsearch/", metricCluster.ClusterID, "/proxy")

	cfg := elasticsearch.Config{
		Addresses: []string{u.String()},
		Username:  ece.user,
		Password:  ece.pass,
		Transport: &tp,
	}

	// Create the Elasticsearch client
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// fmt.Println(es.Cluster.Health())
	doScroll(es, tar)

}

func (ece ECEendpoint) discoverMetricsCluster(tp *CustomRoundTripper) (esCluster, error) {
	client := &http.Client{Transport: tp}
	u, err := url.Parse(ece.eceAPI)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "/api/v1/clusters/elasticsearch")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(ece.user, ece.pass)

	// fmt.Printf("%+v\n", req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	c := new(esClusters)
	json.NewDecoder(resp.Body).Decode(&c)

	for _, cluster := range c.ElasticsearchClusters {
		if cluster.ClusterName == "logging-and-metrics" {
			return cluster, nil
		}
	}
	return esCluster{}, fmt.Errorf("logging-and-metrics not found")
}

type stat struct {
	count      int
	totalBytes int64
	maxSize    int64
}

func doScroll(es *elasticsearch.Client, tar *Tarball) {
	log := logp.NewLogger("MetricScroll")

	query := `{
		"size": 5000,
		"sort": [{"@timestamp": {"order": "desc"}}], 
		"query": {
		  "bool": {
			"filter": {
			  "range": {
				"@timestamp": {"gte": "now-72h"}
			  }
			}
		  }
		}
	  }`

	scrollTimeout := time.Minute
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("*metricbeat*"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithScroll(scrollTimeout),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()

	var scrollQuery esResp
	// unpack json response
	if err = json.NewDecoder(res.Body).Decode(&scrollQuery); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Infof("Total hits from elasticsearch query: %d", scrollQuery.Hits.Total)

	// Create tmp file to write the scroll data into
	file, err := ioutil.TempFile("", "scrollData")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	// gzipWriter := gzip.NewWriter(file)
	gzipWriter, _ := gzip.NewWriterLevel(file, gzip.BestSpeed)
	defer gzipWriter.Close()

	// 50 MiB
	s := stat{maxSize: 52430000}

	err = scrollQuery.writeToFile(gzipWriter, &s)
	if err != nil {
		log.Warn(err)
	}

	scroll := esapi.ScrollRequest{
		ScrollID: scrollQuery.ScrollID,
		Scroll:   scrollTimeout,
	}

	// store the length of hits.hits[]
	var size int
	// loop until size is zero, meaning no more data
	for ok := true; ok; ok = (size > 0) {
		// Perform the request with the client.
		scroller, err := scroll.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer scroller.Body.Close()

		var scrollerResp esResp
		if err = json.NewDecoder(scroller.Body).Decode(&scrollerResp); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// set size
		size = len(scrollerResp.Hits.Hits)

		err = scrollerResp.writeToFile(gzipWriter, &s)
		if err != nil {
			log.Warn(err)
			es.ClearScroll(
				es.ClearScroll.WithScrollID("scrollQuery.ScrollID"),
			)
			break
		}
	}

	es.ClearScroll(
		es.ClearScroll.WithScrollID("scrollQuery.ScrollID"),
	)

	log.Infof("Total fetched from scroll: %d, Total Bytes: %s", s.count, ByteCountBinary(s.totalBytes))

	gzipWriter.Close()
	file.Close()

	fpath, _ := filepath.Abs(file.Name())
	stat, _ := os.Stat(fpath)

	tarRelPath := filepath.Join(cfg.DiagName, "metricbeatData.json.gz")
	tar.AddFile(fpath, stat, tarRelPath)

	defer os.Remove(file.Name())
}

// func (r *esResp) writeToFile(f *os.File) {
func (r *esResp) writeToFile(w io.Writer, s *stat) error {
	buf := &bytes.Buffer{}

	for _, it := range r.Hits.Hits {
		indexAction := []byte(fmt.Sprintf("{\"index\":{\"_index\":\"%s\",\"_type\":\"%s\"}}\n", it.Index, it.Type))
		lineReturn := []byte("\n")
		sourceSize := int64(len(indexAction) + len(it.Source) + len(lineReturn))

		if (s.totalBytes + sourceSize) > s.maxSize {
			w.Write(buf.Bytes())
			return fmt.Errorf("Maximum size reached: %d, limit %d", s.totalBytes, s.maxSize)
		}

		s.count++
		s.totalBytes += sourceSize

		buf.Write(indexAction)
		buf.Write(it.Source)
		buf.Write(lineReturn)
	}

	w.Write(buf.Bytes())
	return nil
}

type CustomRoundTripper struct {
	http.RoundTripper
}

// RoundTrip executes a request and returns a response.
//
func (r CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// req.Header.Set("Accept", "application/yaml")
	// req.Header.Set("X-Request-ID", "foo-123")
	req.Header.Set("X-Management-Request", "true")

	// res, err := http.DefaultTransport.RoundTrip(req)
	resp, err := r.RoundTripper.RoundTrip(req)
	return resp, err
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

type ECEendpoint struct {
	eceAPI string
	user   string
	pass   string
}

type esClusters struct {
	ElasticsearchClusters []esCluster `json:"elasticsearch_clusters"`
}
type esCluster struct {
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

type esResp struct {
	ScrollID string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int64   `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string          `json:"_index"`
			Type   string          `json:"_type"`
			ID     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source json.RawMessage `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
