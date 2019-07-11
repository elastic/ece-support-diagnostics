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

type ECEendpoint struct {
	eceAPI string
	user   string
	pass   string
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
