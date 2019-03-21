package ecediag

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/eluploader/cmd/eluploader/commands"
)

// Config holds configuration variables
type Config struct {
	StartTime     time.Time
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

	RunnerName, err := CheckStoragePath(c.ElasticFolder)
	if err != nil {
		panic(err)
	}
	cfg.RunnerName = RunnerName
	fmt.Println(cfg.RunnerName)

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

// Start is the entry point for the ecediag package
func (c *Config) Start() error {
	cfg = c
	cfg.Initalize()

	fmt.Println(cfg.ElasticFolder)

	l := logp.NewLogger("Main")
	l.Infof("Using %s as temporary storage location", c.Basepath)

	tar := new(Tarball)

	TarFile := filepath.Join(c.Basepath, c.DiagName) + ".tar.gz"
	tar.Create(TarFile)

	defer tar.t.Close()
	defer tar.g.Close()

	runDockerCmds(tar)
	runSystemCmds(tar)
	runCollectLogs(tar)

	// tar.t.Close()
	// tar.g.Close()

	tar.Finalize(filepath.Join(c.Basepath, c.DiagName+".log"))
	runUpload(tar.filepath)

	// add an empty line
	fmt.Println()
	return nil
}

// runUpload is used for the Elastic upload service when the `-u {{ upload uui }}` is present
func runUpload(filePath string) {
	if cfg.UploadUID != "" {
		apiURL := "https://upload-staging.elstc.co"
		numWorkers := runtime.NumCPU()
		cmd := &commands.UploadCmd{UploadID: cfg.UploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
		cmd.Execute()
	}
}
