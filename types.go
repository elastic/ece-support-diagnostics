package ecediag

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
