package eceMetrics

import "encoding/json"

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
