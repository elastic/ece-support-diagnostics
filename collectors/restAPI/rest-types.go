package restAPI

import (
	"encoding/json"
	"net/http"

	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/store"
)

type testFileStore struct {
	store.ContentStore
	cfg *config.Config
}

// Rest defines HTTP URIs to be collected
type Rest struct {
	Method   string
	URI      string
	Filename string
	Headers  map[string]string
	Loop     string
	Sub      []Rest
}

// HTTPClient holds the client, endpoint and credentials
type HTTPClient struct {
	client   *http.Client
	endpoint string
	username string
	passwd   string
	writer   func(filepath string, b []byte) error
}

// RequestTask task
type RequestTask struct {
	config   *HTTPClient
	restItem Rest
}

type v0APIresponse struct {
	Ok           bool   `json:"ok"`
	Message      string `json:"message"`
	EulaAccepted bool   `json:"eula_accepted"`
	Hrefs        struct {
		Regions       string `json:"regions"`
		Elasticsearch string `json:"elasticsearch"`
		Logs          string `json:"logs"`
		DatabaseUsers string `json:"database/users"`
	} `json:"hrefs"`
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
