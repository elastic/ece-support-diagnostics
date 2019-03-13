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
	startTime = time.Now()
	// _, b, _, _  = runtime.Caller(0)
	// Basepath    = filepath.Dir(b)
	// Basepath    = "/tmp"
	Basepath      string
	hostname, _   = os.Hostname()
	ElasticFolder string
	DiagName      string
	TarFile       string
	DisableRest   bool
	UploadUID     string
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
	TarFile = filepath.Join(Basepath, DiagName) + ".tar.gz"
	// tmpfolders := []string{
	// 	filepath.Join(Basepath, "tmp", DiagName, "elastic"),
	// 	filepath.Join(Basepath, "tmp", DiagName, "docker/logs"),
	// }
	// setupFolders(tmpfolders)
}

// Start is an exported function
func Start() error {
	fmt.Println(ElasticFolder)

	l := logp.NewLogger("Main")
	l.Infof("Using %s as temporary storage location", Basepath)

	tar := new(Tarball)
	tar.Create(TarFile)

	defer tar.t.Close()
	defer tar.g.Close()

	runDockerCmds(tar)
	runSystemCmds(tar)

	tar.t.Close()
	tar.g.Close()

	l.Infof("Output: %s", TarFile)
	runUpload(tar.filepath)

	return nil
}

func runUpload(filePath string) {
	if UploadUID != "" {
		apiURL := "https://upload-staging.elstc.co"
		numWorkers := runtime.NumCPU()

		// https://upload-staging.elstc.co/u/f7643a17-8b63-4f17-8a34-85e747c959ea
		// uid := "f7643a17-8b63-4f17-8a34-85e747c959ea"
		cmd := &commands.UploadCmd{UploadID: UploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
		cmd.Execute()
	}
}
