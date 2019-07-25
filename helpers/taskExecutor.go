package helpers

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/elastic/ece-support-diagnostics/store"
)

// type Store interface {
// 	AddData(filePath string, content []byte) error
// }

// type TaskController interface {
// 	TaskCtl
// }

type TaskContext struct {
	Endpoint, Version, User, Pass string
	Meta                          interface{}
	Client                        *http.Client
	Store                         store.ContentStore
	StorePath                     string
	Task
}

// Task is the base unit to execute a web request
//  the attached Meta interface provides the necessary methods
//  and any exported fields used in templating the URL
type Task struct {
	Filename, Method, Uri, Versions string
	RequestBody                     []byte
	Headers                         map[string]string
	Callback                        Callback
}

type Callback interface {
	Exec(taskCtx TaskContext, payload []byte)
}

// Tasks - you know, plural
type Tasks []Task

func (t *Task) httpMethod() string {
	// undefined method field will be treated as "GET"
	if t.Method == "" {
		return http.MethodGet
	}
	return t.Method
}

// URI will return the templated URI string
func (c TaskContext) URI() string {
	return templateString(c.Task.Uri, c.Meta)
}

// URL returns the full URL path to be used in executing the web request
func (c TaskContext) URL() string {
	u, _ := url.Parse(c.Endpoint)
	u.Path = path.Join(u.Path, c.URI())
	url, _ := url.QueryUnescape(u.String())
	return url
}

// Filename is templated and used for writing to the storage adapater
func (c TaskContext) Filename() string {
	return templateString(c.Task.Filename, c.Meta)
}

func (c TaskContext) DoRequest(ctx context.Context) ([]byte, error) {
	// req, _ := http.NewRequest(c.Method(), c.URL(), nil)
	req, _ := http.NewRequest(c.httpMethod(), c.URL(), bytes.NewBuffer(c.RequestBody))
	// Associate the cancellable context
	req = req.WithContext(ctx)

	// Set authentication
	req.SetBasicAuth(c.User, c.Pass)

	// Set headers
	for k, v := range c.Task.Headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		// println("ERROR", err.Error())
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("%s\n", data)
	fpath := filepath.Join(c.StorePath, c.Filename())
	err = c.Store.AddData(fpath, data)
	if err != nil {
		panic(err)
	}
	if c.Callback != nil {
		c.Callback.Exec(c, data)
	}
	return data, nil
}

func (c TaskContext) TaskExecuteWithWaitGroup(wg *sync.WaitGroup) {
	// set a 30s task timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if c.checkVersion() {
		c.DoRequest(ctx)
	} else {
		// println("Skipping due to version check - ", c.Task.versions)
	}
	wg.Done()
}

func (c TaskContext) TaskExecute() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if c.checkVersion() {
		c.DoRequest(ctx)
	} else {
		// println("Skipping due to version check - ", c.Task.versions)
	}
}

func (c TaskContext) checkVersion() bool {
	// If no version is defined, it will always be executed
	if c.Task.Versions == "" {
		return true
	}
	// Check the version against the constraint
	chk, _ := semver.NewConstraint(c.Task.Versions)
	v, _ := semver.NewVersion(c.Version)
	return chk.Check(v)
}

func templateString(item string, Obj interface{}) string {
	t := template.Must(template.New("testing").Parse(item))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, Obj)

	// Need to figure out how to validate the data upon build
	// All built in data should properly parse / template.
	if err != nil {
		// println(err.Error())
	}
	return tpl.String()
}
