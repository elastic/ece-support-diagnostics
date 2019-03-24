package ecediag

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/elastic/beats/libbeat/logp"
	"golang.org/x/crypto/ssh/terminal"
)

// HTTPClient holds the client, endpoint and credentials
type HTTPClient struct {
	client   *http.Client
	endpoint string
	username string
	passwd   string
}

// RequestTask task
type RequestTask struct {
	config   *HTTPClient
	restItem Rest
}

// RunRest starts the chain of functions to collect the Rest/HTTP calls
func RunRest(d types.Container, tar *Tarball) {
	// var err error

	httpClient := NewClient()
	httpClient.endpoint = "https://0.0.0.0:12443/"
	err := httpClient.SetupCredentials()
	panicError(err)

	fmt.Println("[ ] Collecting API information ECE and Elasticsearch")

	var wg sync.WaitGroup
	for _, item := range rest {
		wg.Add(1)
		task := RequestTask{config: httpClient, restItem: item}
		go task.fetch(tar, &wg)
	}
	wg.Wait()

	clearStdoutLine()
	fmt.Println("[✔] Collected API information ECE and Elasticsearch")
}

// NewClient returns the configured http client to be used for http requests
func NewClient() *HTTPClient {
	var tr = &http.Transport{
		// Disable Certificate Checking
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Connection timeout = 5s
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		// TLS Handshake Timeout = 5s
		TLSHandshakeTimeout: 5 * time.Second,
	}
	// HTTP Timeout = 10s
	myClient := &http.Client{Timeout: 10 * time.Second, Transport: tr}
	return &HTTPClient{client: myClient}
}

// SetupCredentials checks that the auth credentials are valid
//  successful auth creds are used for remaining requests
func (r *HTTPClient) SetupCredentials() error {
	log := logp.NewLogger("ValidateAuth")

	r.username, r.passwd = getCredentials()

	req, err := http.NewRequest("GET", r.endpoint+"api/v0", nil)
	if err != nil {
		// handle err
	}

	fmt.Println()

	req.SetBasicAuth(r.username, r.passwd)
	resp, err := r.client.Do(req)
	panicError(err)

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		defer resp.Body.Close()
		v0Response := new(v0APIresponse)
		json.NewDecoder(resp.Body).Decode(v0Response)

		if v0Response.Ok {
			for i := 0; i <= 2; i++ {
				clearStdoutLine()
			}
			fmt.Printf("Authenticated\n")
			fmt.Printf("\t✔ Username (%s)\n", r.username)
			fmt.Printf("\t✔ Password\n")

			log.Infof("Cloud UI Resolved, using %s", req.URL)
			return nil
		}
	}
	return fmt.Errorf("Authentication failed")
}

// fetch dispatches the Rest/HTTP request
func (parent *RequestTask) fetch(tar *Tarball, wg *sync.WaitGroup) {

	url := parent.config.endpoint + strings.TrimLeft(parent.restItem.URI, "/")
	req, err := http.NewRequest("GET", url, nil)
	// TODO: handle err?

	req.SetBasicAuth(parent.config.username, parent.config.passwd)
	req.Header.Set("X-Management-Request", "true")
	resp, err := parent.config.client.Do(req)
	if err != nil {
		log.Fatal(url, err)
	}
	// if resp.StatusCode >= 200 && resp.StatusCode <= 299 { }
	respBody, err := ioutil.ReadAll(resp.Body)

	archiveFile := filepath.Join(cfg.DiagName, parent.restItem.Filename)
	tar.AddData(archiveFile, respBody)

	parent.checkSubItems(respBody, tar)

	wg.Done()
}

// checkSubItems is used when `Sub` is defined in the Rest object, and contains a `Loop` item.
//  It tries to unpack the parent JSON response into a map, and assert the proper type (array/object)
func (parent *RequestTask) checkSubItems(respBody []byte, tar *Tarball) {

	if len(parent.restItem.Sub) > 0 {

		// TODO, json decode error should fail.
		var resp interface{}
		resp = readJSON(respBody)

		switch json := resp.(type) {

		// Json response for the parent Rest response is a JSON Array
		case []interface{}:
			// TODO, this should only happen when loop is specified.
			fmt.Println(json, "Array!")
			parent.subLoop(json, tar)
			// parent.iterateSub(resp, tar)

		// Json response for parent Rest response is a JSON Object
		case map[string]interface{}:
			if parent.restItem.Loop == "" {
				// Iterate with top level map
				parent.iterateSub(resp, tar)
			} else {
				if val, ok := json[parent.restItem.Loop]; ok {
					switch data := val.(type) {
					case []interface{}:
						parent.subLoop(data, tar)
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

func (parent *RequestTask) subLoop(respArray []interface{}, tar *Tarball) {
	for _, Item := range respArray {
		parent.iterateSub(Item, tar)
	}
}

func (parent *RequestTask) iterateSub(It interface{}, tar *Tarball) {
	var wg sync.WaitGroup
	l := logp.NewLogger("Elasticsearch")

	s := It.(map[string]interface{})
	l.Infof("Gathering cluster diagnostic: %s, %s", s["cluster_id"], s["cluster_name"])

	for _, item := range parent.restItem.Sub {
		wg.Add(1)
		// render template
		item.templater(It)
		task := RequestTask{config: parent.config, restItem: item}
		go task.fetch(tar, &wg)
	}
	wg.Wait()
}

// templater when called runs templating for the defined fields
func (R *Rest) templater(Obj interface{}) {
	R.Filename = runTemplate(R.Filename, Obj)
	R.URI = runTemplate(R.URI, Obj)
}

// runTemplate performs the string substitution using the html/template package
func runTemplate(item string, Obj interface{}) string {
	t := template.Must(template.New("testing").Parse(item))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, Obj)
	if err != nil {
		log.Println("executing template:", Obj)
	}
	return tpl.String()
}

// readJSON unpacks the Rest/HTTP request into a generic interface
func readJSON(in []byte) interface{} {
	var data interface{}
	err := json.Unmarshal(in, &data)
	panicError(err)
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
	panicError(err)
	// if err == nil {fmt.Println("\nPassword typed: " + string(bytePassword))}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password)
	// return "readonly", strings.TrimSpace(password)
}
