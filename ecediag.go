package ecediag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/eluploader/cmd/eluploader/commands"
)

const ()

var (
	startTime   = time.Now()
	hostname, _ = os.Hostname()
	// _, b, _, _  = runtime.Caller(0)

	// Basepath provides the tmp location to create the tar file in
	Basepath string

	// ElasticFolder provides the path to where ECE is installed
	ElasticFolder string

	// DiagName is used for the output tar file name
	DiagName string

	// DisableRest is used for disabling collecting Rest/HTTP requests
	DisableRest bool

	// UploadUID provides the unique ID that needs to be specified for using the Elastic upload service
	UploadUID string
)

func init() {
	flag.StringVar(&ElasticFolder, "f", "/mnt/data/elastic", "Path to the elastic folder")
	flag.StringVar(&Basepath, "t", "/tmp", "Path to the elastic folder")
	flag.StringVar(&UploadUID, "u", "", "Elastic Upload ID")
	flag.BoolVar(&DisableRest, "disableRest", false, "Disable Rest calls")
	flag.Parse()
	RunnerName, err := checkStoragePath()
	if err != nil {
		panic(err)
	}
	DiagDate := fmt.Sprintf("-%d%02d%02d-%02d%02d%02d",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())
	DiagName = "ecediag-" + RunnerName + DiagDate

	config := logp.Config{
		Beat:       "ece-support-diag",
		JSON:       false,
		Level:      logp.DebugLevel,
		ToStderr:   false,
		ToSyslog:   false,
		ToFiles:    true,
		ToEventLog: false,
		Files: logp.FileConfig{
			Path:        Basepath,
			Name:        DiagName + ".log",
			MaxSize:     20000000,
			MaxBackups:  0,
			Permissions: 0644,
			// Interval:       4 * time.Hour,
			RedirectStderr: true,
		},
	}
	logp.Configure(config)

	// tmpfolders := []string{
	// 	filepath.Join(Basepath, "tmp", DiagName, "elastic"),
	// 	filepath.Join(Basepath, "tmp", DiagName, "docker/logs"),
	// }
	// setupFolders(tmpfolders)
}

// Start is the entry point for the ecediag package
func Start() error {
	fmt.Println(ElasticFolder)

	l := logp.NewLogger("Main")
	l.Infof("Using %s as temporary storage location", Basepath)

	tar := new(Tarball)

	TarFile := filepath.Join(Basepath, DiagName) + ".tar.gz"
	tar.Create(TarFile)

	defer tar.t.Close()
	defer tar.g.Close()

	runDockerCmds(tar)
	runSystemCmds(tar)

	// tar.t.Close()
	// tar.g.Close()

	tar.Finalize(filepath.Join(Basepath, DiagName+".log"))
	runUpload(tar.filepath)

	return nil
}

// runUpload is used for the Elastic upload service when the `-u {{ upload uui }}` is present
func runUpload(filePath string) {
	if UploadUID != "" {
		apiURL := "https://upload-staging.elstc.co"
		numWorkers := runtime.NumCPU()
		cmd := &commands.UploadCmd{UploadID: UploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
		cmd.Execute()
	}
}
