package config

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/elastic/ece-support-diagnostics/discovery"
	"github.com/elastic/ece-support-diagnostics/store"
)

// Config holds configuration variables
type Config struct {
	// internal
	StartTime time.Time

	// args set default
	OlderThan     time.Duration
	Basepath      string
	ElasticFolder string
	DisableRest   bool
	UploadUID     string

	Store store.ContentStore

	// Runner name needs to be discovered. Needs to be init
	diagName   string
	runnerName string

	// Needs to be discovered
	APIendpoint string
	Auth

	// automatically created
	HTTPclient *http.Client
}

// New returns a test config
func New() *Config {
	return &Config{
		StartTime:  time.Now(),
		Basepath:   "",
		HTTPclient: NewHTTPClient(),
	}
}

// Initalize makes sure runtime variables are all set
func (c *Config) Initalize() {
	RunnerName, err := discovery.CheckStoragePath(c.ElasticFolder)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(`
	please set the -f, --eceHomeFolder argument to the appropriate 
	folder path where ECE is installed, or check that the user 
	running this diagnostic utility has the appropriate permissions 
	to that location
		`)
		os.Exit(1)
	}
	c.runnerName = RunnerName
	// fmt.Printf("ECE Runner Name: %s\n", c.runnerName)

	c.diagName = "ecediag-" + c.runnerName + c.diagDate()

	c.setupLogging()

	// need logging from this point
	if c.DisableRest != true {
		c.APIendpoint, err = discovery.DiscoverAPI(c.ElasticFolder, c.HTTPclient)
		if err != nil {
			fmt.Println("Could not determine coordinator API endpoint to use")
			fmt.Println("\u26A0 \u26A0 \u26A0 API support data will be skipped \u26A0 \u26A0 \u26A0")

			// no endpoint to make rest calls to, skip at this point
			c.DisableRest = true
			return
		}
		// fmt.Printf(helpers.PreviousLine)
		fmt.Println("Using", c.APIendpoint, "for API calls")

		err := c.initalizeCredentials()
		if err != nil {
			// fmt.Printf("\n\u26A0  \u26A0  \u26A0    %s    \u26A0  \u26A0  \u26A0\n", err.Error())
			fmt.Printf("\n⚠  ⚠  ⚠  ⚠  ⚠    %s    ⚠  ⚠  ⚠  ⚠  ⚠\n", err.Error())
			fmt.Printf(`
/******************************************************/
/* If you are having trouble authenticating, you can  */ 
/* skip collection of the API based support data by   */
/* using the --disableRest flag. Please understand    */
/* this will severely limit Elastic Support's ability */
/* to provide timely help.                            */
/******************************************************/

`)

			// TODO: Need to safely exit and cleanup files
			os.Exit(1)
		}
	}
}

// RunnerName runnerName
func (c *Config) RunnerName() string {
	return c.runnerName
}

// DiagnosticFilename will be the output filename without any extension appended
func (c *Config) DiagnosticFilename() string {
	// make sure we have initalized
	if c.runnerName == "" {
		println(c.runnerName)
		panic("DiagnosticFilename() has not been initalized")
	}
	return c.diagName
}

// DiagnosticTarFilePath provides the full filepath to the destination tar file
func (c *Config) DiagnosticTarFilePath() string {
	return filepath.Join(c.Basepath, c.DiagnosticFilename()+".tar.gz")
}

// DiagnosticLogFilePath provides the full filepath to the destination log file
func (c *Config) DiagnosticLogFilePath() string {
	return filepath.Join(c.Basepath, c.DiagnosticFilename()+".log")
}

func (c *Config) diagDate() string {
	return fmt.Sprintf("-%d%02d%02d-%02d%02d%02d",
		c.StartTime.Year(),
		c.StartTime.Month(),
		c.StartTime.Day(),
		c.StartTime.Hour(),
		c.StartTime.Minute(),
		c.StartTime.Second(),
	)
}
