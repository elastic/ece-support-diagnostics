package ece_support_diagnostic

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

// Hack for testing only.
// TODO: Remove this and setup a proper storage path.
var (
	_, b, _, _   = runtime.Caller(0)
	Basepath     = filepath.Dir(b)
	hostname, _  = os.Hostname()
	DiagName     = "ece_diag_" + hostname + time.Now().Format("_2006-01-02_15:04:05_MST")
	SystemTmpDir = filepath.Join(Basepath, "tmp", DiagName, "elastic")
	DockerTmpDir = filepath.Join(Basepath, "tmp", DiagName, "docker")
)

// Run is an exported function
func Run() error {
	logp.Info("Hello world")
	logp.Info(Basepath)
	setupFolders()
	SystemCommands()
	DockerCommands()
	TarLogs()
	return nil
}

func setupFolders() {
	tmp_folders := []string{
		filepath.Join(Basepath, "tmp", DiagName, "elastic"),
		filepath.Join(Basepath, "tmp", DiagName, "docker/logs"),
	}
	for _, folder := range tmp_folders {
		f, _ := filepath.Abs(folder)
		flog := logp.NewLogger("folders")
		flog.Info("Temp folder setup: ", f)
		os.MkdirAll(f, os.ModePerm)
	}
}
