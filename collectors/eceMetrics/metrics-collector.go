package eceMetrics

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
	elasticsearch "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/tidwall/gjson"
)

type metricsCollector struct {
	log *logp.Logger
	tp  CustomRoundTripper
}

// CustomRoundTripper is used to override the http.RoundTripper and attach headers
//  to each request
type CustomRoundTripper struct {
	http.RoundTripper
}

// Run runs
func Run(status chan<- string, cfg *config.Config) {

	if cfg.DisableRest == true {
		status <- fmt.Sprintf("\u26A0 skipping collection of ECE metricbeat data")
		return
	}

	metrics := metricsCollector{
		log: logp.NewLogger("MetricScroll"),
		tp: CustomRoundTripper{
			RoundTripper: cfg.HTTPclient.Transport,
		},
	}

	err := metrics.ScrollRunner(cfg)
	if err != nil {
		status <- fmt.Sprintf("\u2715 Failed to collect ECE metricbeat data")
	}

	status <- fmt.Sprintf("\u2713 collected ECE metricbeat data")
}

func (m metricsCollector) ScrollRunner(cfg *config.Config) error {
	metricCluster, err := m.discoverMetricsCluster(cfg)
	if err != nil {
		m.log.Error(err)
		return err
	}
	u, err := url.Parse(cfg.APIendpoint)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "/api/v1/clusters/elasticsearch/", metricCluster, "/proxy")

	esConf := elasticsearch.Config{
		Addresses: []string{u.String()},
		Username:  cfg.Auth.User,
		Password:  cfg.Auth.Pass,
		Transport: &m.tp,
	}

	// Create the Elasticsearch client
	es, err := elasticsearch.NewClient(esConf)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// fmt.Println(es.Cluster.Health())
	m.doScroll(es, cfg)
	return nil
}

func (m metricsCollector) discoverMetricsCluster(cfg *config.Config) (string, error) {
	client := &http.Client{Transport: &m.tp}

	if cfg.APIendpoint == "" {
		return "", fmt.Errorf("API endpoint is blank")
	}
	u, _ := url.Parse(cfg.APIendpoint)
	u.Path = path.Join(u.Path, "/api/v1/clusters/elasticsearch")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(cfg.Auth.User, cfg.Auth.Pass)

	// fmt.Printf("%+v\n", req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// TODO: check for error?
	respBytes, _ := ioutil.ReadAll(resp.Body)

	// select the logging-and-metrics cluster
	clusterID := gjson.GetBytes(respBytes, `elasticsearch_clusters.#(cluster_name=="logging-and-metrics").cluster_id`)
	if !clusterID.Exists() {
		return "", fmt.Errorf("Could not find logging-and-metrics cluster ID")
	}
	return clusterID.String(), nil
}

type stat struct {
	count      int
	totalBytes int64
	maxSize    int64
}

func (m metricsCollector) doScroll(es *elasticsearch.Client, cfg *config.Config) {
	log := logp.NewLogger("MetricScroll")

	// TODO: make filter for event.dataset: `system.process` (keyword) optional

	query := `{
		"size": 5000,
		"sort": [
		  { "@timestamp": { "order": "desc" }}
		],
		"query": {
		  "bool": {
			"filter": [
			  { "range": { "@timestamp": { "gte": "now-72h" }}},
			  {
				"bool": {
				  "must_not": [
					{
					  "term": {
						"event.dataset": {
						  "value": "system.process"
						}
					  }
					}
				  ]
				}
			  }
			]
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

	log.Infof("Total fetched from scroll: %d, Total Bytes: %s", s.count, helpers.ByteCountBinary(s.totalBytes))

	gzipWriter.Close()
	file.Close()

	fpath, _ := filepath.Abs(file.Name())
	stat, _ := os.Stat(fpath)

	tarRelPath := filepath.Join(cfg.DiagnosticFilename(), "metricbeatData.json.gz")
	cfg.Store.AddFile(fpath, stat, tarRelPath)

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
