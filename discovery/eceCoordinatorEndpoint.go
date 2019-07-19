package discovery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/tidwall/gjson"
)

// DiscoverAPI determines the appropriate ECE coordinator endpoint to use
func DiscoverAPI(installFolder string, r *http.Client) (string, error) {
	// ./*/services/client-forwarder/managed/ensemble-state-file.json

	clientForwarderEnsemble := "*/services/client-forwarder/managed/ensemble-state-file.json"
	globFolder := filepath.Join(installFolder, clientForwarderEnsemble)

	// DEBUG: fmt.Println(globFolder)

	matches, err := filepath.Glob(globFolder)
	if err != nil {
		panic(err)
	}

	if len(matches) > 0 {
		// TODO: there never should be more than one match, handle error if more than 1?
		ensembleFile := matches[0]
		ensemble, _ := ioutil.ReadFile(ensembleFile)

		result := gjson.GetBytes(ensemble, "coordinators")

		// TODO: should only need one transport/http client. Config?
		// r := config.NewHttpClient()

		var endpoint string
		result.ForEach(func(key, value gjson.Result) bool {
			// Try each Coordinator host to see if the API is reachable

			publicHostname := value.Get("public_hostname").String()

			// DEBUG: println(publicHostname)

			println("Checking connection:", publicHostname)
			// Check if :12443/api/v1 can be reached
			endpoint, err = checkEndpoint(r, publicHostname)
			if err != nil {
				// endpoint did not work
				println(err.Error())
				return true // keep iterating
			}

			// valid endpoint, stop iteration
			return false

		})

		if endpoint != "" {
			fmt.Println(endpoint)
			return endpoint, nil
		}
		// fmt.Println(err)
		return "", err
	}

	// Couldn't find the client-forwarder ensemble-state-file.json
	return "", fmt.Errorf("Couldn't find the client-forwarder ensemble-state-file.json")
}

func checkEndpoint(r *http.Client, host string) (string, error) {
	endpoint := "https://" + host + ":12443"
	u, _ := url.Parse(endpoint)
	u.Path = path.Join(u.Path, "api/v1")

	// fmt.Println(u.String())
	req, _ := http.NewRequest("GET", u.String(), nil)

	resp, err := r.Do(req)
	if err != nil {
		// FAILED to connect
		return "", err
	}
	if resp.StatusCode == 401 {
		// HTTP/1.1 401 Unauthorized
		// Looking for unauthorized response to know that we can talk to the endpoint
		return endpoint, nil
	}
	return "", fmt.Errorf("Endpoint not valid")
}
