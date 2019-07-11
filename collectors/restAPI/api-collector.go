package restAPI

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/discovery"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"github.com/elastic/ece-support-diagnostics/store"
	"golang.org/x/crypto/ssh/terminal"
)

func Run(t store.ContentStore, rest []Rest, config *config.Config) {
	store := testFileStore{t, config}
	store.runRest(rest)
}

// NewClient returns the configured http client to be used for http requests
//  TODO: look into consolidating to single transport/http client.
func newClient() *HTTPClient {
	return &HTTPClient{client: config.NewHttpClient()}
}

// RunRest starts the chain of functions to collect the Rest/HTTP calls
func (t testFileStore) runRest(rest []Rest) {
	httpClient := newClient()
	httpClient.writer = t.AddData

	// TODO: create a safe method to return the endpoint + URI.
	//  when switching to the discovery, the lack of the trailing slash caused some issues
	// httpClient.endpoint = "https://0.0.0.0:12443/"
	endpoint, err := discovery.DiscoverAPI(t.cfg.ElasticFolder, httpClient.client)
	if err != nil {
		fmt.Println("Failed to discover ECE API to collect data from")
		return
	}
	httpClient.endpoint = endpoint

	err = t.setupCredentials(httpClient)
	helpers.PanicError(err)

	fmt.Println("[ ] Collecting ECE metricbeat data")
	creds := &ECEendpoint{
		eceAPI: httpClient.endpoint,
		user:   httpClient.username,
		pass:   httpClient.passwd,
	}
	t.ScrollRunner(creds)
	helpers.ClearStdoutLine()
	fmt.Println("[✔] Collected ECE metricbeat data")

	fmt.Println("[ ] Collecting API information ECE and Elasticsearch")

	var wg sync.WaitGroup
	for _, item := range rest {
		wg.Add(1)
		task := RequestTask{config: httpClient, restItem: item}
		go t.fetch(&task, &wg)
	}
	wg.Wait()

	helpers.ClearStdoutLine()
	fmt.Println("[✔] Collected API information ECE and Elasticsearch")
}

// SetupCredentials checks that the auth credentials are valid
//  successful auth creds are used for remaining requests
func (t testFileStore) setupCredentials(r *HTTPClient) error {
	log := logp.NewLogger("ValidateAuth")

	r.username, r.passwd = getCredentials()

	url, _ := url.Parse(r.endpoint)
	// TODO: check that this API endpoint will work on 1.x versions?
	url.Path = path.Join(url.Path, "/api/v1/platform/license")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		// handle err?
		// need to understand if there is any risk of an error in creating the http request?
	}

	fmt.Println()

	req.SetBasicAuth(r.username, r.passwd)
	resp, err := r.client.Do(req)

	// auth failed? retry?
	helpers.PanicError(err)

	if resp.StatusCode == 200 || resp.StatusCode == 400 {
		for i := 0; i <= 2; i++ {
			helpers.ClearStdoutLine()
		}
		fmt.Printf("Authenticated\n")
		fmt.Printf("\t✔ Username (%s)\n", r.username)
		fmt.Printf("\t✔ Password\n")

		log.Infof("Cloud UI Resolved, using %s", req.URL)
		return nil
	}

	// TODO: write license response to file?

	return fmt.Errorf("Authentication failed")
}

// fetch dispatches the Rest/HTTP request
func (t testFileStore) fetch(parent *RequestTask, wg *sync.WaitGroup) {

	url, err := url.Parse(parent.config.endpoint)
	url.Path = path.Join(url.Path, parent.restItem.URI)
	// url := parent.config.endpoint + strings.TrimLeft(parent.restItem.URI, "/")

	req, err := http.NewRequest("GET", url.String(), nil)
	// TODO: handle err?

	// set auth
	req.SetBasicAuth(parent.config.username, parent.config.passwd)
	// set headers
	for k, v := range parent.restItem.Headers {
		req.Header.Set(k, v)
	}
	// req.Header.Set("X-Management-Request", "true")
	resp, err := parent.config.client.Do(req)
	if err != nil {
		log.Fatal(url, err)
	}
	// if resp.StatusCode >= 200 && resp.StatusCode <= 299 { }
	respBody, err := ioutil.ReadAll(resp.Body)

	archiveFile := filepath.Join(t.cfg.DiagName, parent.restItem.Filename)

	// write response data to file
	err = parent.config.writer(archiveFile, respBody)
	helpers.PanicError(err)

	t.checkSubItems(parent, respBody)

	wg.Done()
}

// checkSubItems is used when `Sub` is defined in the Rest object, and contains a `Loop` item.
//  It tries to unpack the parent JSON response into a map, and assert the proper type (array/object)
func (t testFileStore) checkSubItems(parent *RequestTask, respBody []byte) {

	if len(parent.restItem.Sub) > 0 {

		// TODO, json decode error should fail.
		var resp interface{}
		resp = readJSON(respBody)

		switch json := resp.(type) {

		// Json response for the parent Rest response is a JSON Array
		case []interface{}:
			// TODO, this should only happen when loop is specified.
			fmt.Println(json, "Array!")
			t.subLoop(parent, json)
			// parent.iterateSub(resp)

		// Json response for parent Rest response is a JSON Object
		case map[string]interface{}:
			if parent.restItem.Loop == "" {
				// Iterate with top level map
				t.iterateSub(parent, resp)
			} else {
				if val, ok := json[parent.restItem.Loop]; ok {
					switch data := val.(type) {
					case []interface{}:
						t.subLoop(parent, data)
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

func (t testFileStore) subLoop(parent *RequestTask, respArray []interface{}) {
	for _, Item := range respArray {
		t.iterateSub(parent, Item)
	}
}

func (t testFileStore) iterateSub(parent *RequestTask, It interface{}) {
	var wg sync.WaitGroup
	l := logp.NewLogger("Elasticsearch")

	s := It.(map[string]interface{})
	l.Infof("Gathering cluster diagnostic: %s, %s", s["cluster_id"], s["cluster_name"])

	for _, item := range parent.restItem.Sub {
		wg.Add(1)
		// render template
		item.templateService(It)
		task := RequestTask{config: parent.config, restItem: item}
		go t.fetch(&task, &wg)
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

// getCredentials is used for securely prompting for a password from stdin
//  it uses the x/crypto/ssh/terminal package to ensure stdin echo is disabled
func getCredentials() (usr, pass string) {
	fmt.Println("Please Enter Your ECE Admin Credentials")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	// fmt.Println("Username (read-only)")
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	helpers.PanicError(err)
	// if err == nil {fmt.Println("\nPassword typed: " + string(bytePassword))}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password)
	// return "readonly", strings.TrimSpace(password)
}
