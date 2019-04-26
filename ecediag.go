package ecediag

import (
	"fmt"
	"path/filepath"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/collectors/docker"
	"github.com/elastic/ece-support-diagnostics/collectors/restAPI"
	"github.com/elastic/ece-support-diagnostics/collectors/systemInfo"
	"github.com/elastic/ece-support-diagnostics/collectors/systemLogs"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/store/tar"
	"github.com/elastic/ece-support-diagnostics/uploader"
)

// Start is the entry point for the ecediag package
func Start(c *config.Config) error {

	cfg := c
	cfg.Initalize()

	fmt.Printf("Using %s for ECE install location\n", cfg.ElasticFolder)

	l := logp.NewLogger("Main")
	l.Infof("Using %s as temporary storage location", c.Basepath)

	// tar := new(Tarball)

	tarFilepath := filepath.Join(c.Basepath, c.DiagName) + ".tar.gz"

	tar, err := tar.Create(tarFilepath)
	defer tar.Close()

	if err != nil {
		// handle err
	}

	// runDockerCmds(tar)
	docker.Run(tar, cfg)
	restAPI.Run(tar, rest, cfg)
	systemInfo.Run(tar, SystemCmds, SystemFiles, cfg)
	systemLogs.Run(tar, cfg)

	logfilePath := filepath.Join(c.Basepath, c.DiagName+".log")
	logTarPath := filepath.Join(cfg.DiagName, "diagnostic.log")
	tar.Finalize(logfilePath, logTarPath)

	if cfg.UploadUID != "" {
		uploader.RunUpload(tar.Filepath(), cfg)
	}
	// add an empty line
	fmt.Println()
	return nil

}
