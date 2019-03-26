package ecediag

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

// Config holds configuration variables
type Config struct {
	StartTime     time.Time
	OlderThan     time.Duration
	Basepath      string
	ElasticFolder string
	DiagName      string
	RunnerName    string
	DisableRest   bool
	UploadUID     string
}

var cfg = New()

// New returns a test config
func New() *Config {
	c := &Config{
		StartTime: time.Now(),
		Basepath:  "",
	}
	return c
}

// Initalize makes sure runtime variables are all set
func (c *Config) Initalize() {
	cfg.OlderThan = 72 * time.Hour
	RunnerName, err := CheckStoragePath(c.ElasticFolder)
	if err != nil {
		panic(err)
	}
	cfg.RunnerName = RunnerName
	fmt.Printf("Runner Name: %s\n", cfg.RunnerName)

	DiagDate := fmt.Sprintf("-%d%02d%02d-%02d%02d%02d",
		c.StartTime.Year(),
		c.StartTime.Month(),
		c.StartTime.Day(),
		c.StartTime.Hour(),
		c.StartTime.Minute(),
		c.StartTime.Second(),
	)
	c.DiagName = "ecediag-" + RunnerName + DiagDate

	config := logp.Config{
		Beat:       "ece-support-diag",
		JSON:       false,
		Level:      logp.DebugLevel,
		ToStderr:   false,
		ToSyslog:   false,
		ToFiles:    true,
		ToEventLog: false,
		Files: logp.FileConfig{
			Path:        c.Basepath,
			Name:        c.DiagName + ".log",
			MaxSize:     20000000,
			MaxBackups:  0,
			Permissions: 0644,
			// Interval:       4 * time.Hour,
			RedirectStderr: true,
		},
	}
	logp.Configure(config)

}
