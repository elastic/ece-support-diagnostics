package eceAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"sync"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/helpers"
)

// Run is the entry point to make the ECE Rest API calls
func Run(rest []Rest, cfg *config.Config) {
	fmt.Println("[ ] Collecting API information ECE and Elasticsearch")

	var wg sync.WaitGroup
	for _, item := range rest {
		wg.Add(1)
		task := item
		go task.DoRequest(cfg, &wg)
	}
	wg.Wait()

	helpers.ClearStdoutLine()
	fmt.Println("[✔] Collected API information ECE and Elasticsearch")
}

// DoRequest dispatches the Rest/HTTP request
func (task Rest) DoRequest(cfg *config.Config, wg *sync.WaitGroup) {

	url, err := url.Parse(cfg.APIendpoint)
	url.Path = path.Join(url.Path, task.URI)

	req, err := http.NewRequest("GET", url.String(), nil)
	// TODO: handle err?

	// set auth
	req.SetBasicAuth(cfg.Auth.User, cfg.Auth.Pass)
	// set headers
	for k, v := range task.Headers {
		req.Header.Set(k, v)
	}
	// req.Header.Set("X-Management-Request", "true")

	resp, err := cfg.HTTPclient.Do(req)
	if err != nil {
		log.Fatal(url, err)
	}
	// if resp.StatusCode >= 200 && resp.StatusCode <= 299 { }
	respBody, err := ioutil.ReadAll(resp.Body)

	archiveFile := filepath.Join(cfg.DiagnosticFilename(), task.Filename)

	// write response data to file
	err = cfg.Store.AddData(archiveFile, respBody)
	helpers.PanicError(err)

	task.checkSubItems(cfg, respBody)

	wg.Done()
}

// checkSubItems is used when `Sub` is defined in the Rest object, and contains a `Loop` item.
//  It tries to unpack the parent JSON response into a map, and assert the proper type (array/object)
// func (t testFileStore) checkSubItems(parent *RequestTask, respBody []byte) {
func (task Rest) checkSubItems(cfg *config.Config, respBody []byte) {

	if len(task.Sub) > 0 {

		// TODO, json decode error should fail.
		var resp interface{}
		resp = readJSON(respBody)

		switch json := resp.(type) {

		// Json response for the parent Rest response is a JSON Array
		case []interface{}:
			// TODO, this should only happen when loop is specified.
			fmt.Println(json, "Array!")
			task.subLoop(cfg, json)
			// parent.iterateSub(resp)

		// Json response for parent Rest response is a JSON Object
		case map[string]interface{}:
			if task.Loop == "" {
				// Iterate with top level map
				task.iterateSub(cfg, resp)
			} else {
				if val, ok := json[task.Loop]; ok {
					switch data := val.(type) {
					case []interface{}:
						task.subLoop(cfg, data)
					default:
						// break, the key specified is not an array
					}
				} else {
					// key error, the specified key did not exist
				}
			}

		default:
			fmt.Println("SHIT!!!!")
		}
	}
}

func (task Rest) subLoop(cfg *config.Config, respArray []interface{}) {
	for _, Item := range respArray {
		task.iterateSub(cfg, Item)
	}
}

func (task Rest) iterateSub(cfg *config.Config, data interface{}) {
	var wg sync.WaitGroup
	l := logp.NewLogger("Elasticsearch")

	s := data.(map[string]interface{})
	l.Infof("Gathering cluster diagnostic: %s, %s", s["cluster_id"], s["cluster_name"])

	for _, item := range task.Sub {
		wg.Add(1)
		// render template
		item.templateService(data)
		subTask := item
		go subTask.DoRequest(cfg, &wg)
	}
	wg.Wait()
}

// readJSON unpacks the Rest/HTTP request into a generic interface
func readJSON(in []byte) interface{} {
	var data interface{}
	err := json.Unmarshal(in, &data)
	helpers.PanicError(err)
	return data
}
