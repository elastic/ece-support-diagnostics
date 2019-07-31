package config

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/helpers"
	"golang.org/x/crypto/ssh/terminal"
)

// Auth holds the user and pass
type Auth struct {
	User, Pass    string
	authenticated bool
}

// initalizeCredentials ensures that auth credentials are setup and valid
func (c *Config) initalizeCredentials() error {
	log := logp.NewLogger("ValidateAuth")

	c.checkForPassword()

	url, _ := url.Parse(c.APIendpoint)

	// TODO: check that this API endpoint will work on 1.x versions?
	url.Path = path.Join(url.Path, "/api/v1/platform/license")

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		// handle err?
		// need to understand if there is any risk of an error in creating the http request?
	}

	fmt.Println()

	req.SetBasicAuth(c.Auth.User, c.Auth.Pass)
	resp, err := c.HTTPclient.Do(req)

	// auth failed? retry?
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 400 {
		for i := 0; i <= 2; i++ {
			helpers.ClearStdoutLine()
		}
		fmt.Printf("Authenticated\n")
		fmt.Printf("\t✔ Username (%s)\n", c.Auth.User)
		fmt.Printf("\t✔ Password\n")

		log.Infof("Cloud UI Resolved, using %s", req.URL)
		return nil
	}

	// TODO: write license response to file?

	return fmt.Errorf("Authentication failed")
}

func (c *Config) checkForPassword() {
	if c.Auth.User != "" {
		// Split the username if it contains a `:`
		auth := strings.SplitN(c.Auth.User, ":", 2)
		if len(auth) == 2 {
			c.Auth.User = auth[0]
			c.Auth.Pass = auth[1]
		} else {
			// Only the username was specified
			c.Auth.Pass = promptForPassword()
		}
	} else {
		// no username flag invoked, need user & pass
		c.Auth.User, c.Auth.Pass = credsFromCmdPrompt()
	}
}

func promptForPassword() string {
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	password := string(bytePassword)
	return password
}

// getCredentials is used for securely prompting for a password from stdin
//  it uses the x/crypto/ssh/terminal package to ensure stdin echo is disabled
func credsFromCmdPrompt() (usr, pass string) {
	fmt.Println("Please Enter Your ECE Admin Credentials")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	// fmt.Println("Username (read-only)")
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	// if err == nil {fmt.Println("\nPassword typed: " + string(bytePassword))}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password)
	// return "readonly", strings.TrimSpace(password)
}
