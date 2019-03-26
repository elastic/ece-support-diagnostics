package ecediag

import (
	"fmt"
	"path/filepath"

	"github.com/elastic/beats/libbeat/logp"
)

// Start is the entry point for the ecediag package
func (c *Config) Start() error {
	cfg = c
	cfg.Initalize()

	fmt.Printf("Using %s for ECE install location\n", cfg.ElasticFolder)

	l := logp.NewLogger("Main")
	l.Infof("Using %s as temporary storage location", c.Basepath)

	// tar := new(Tarball)

	tarFile := filepath.Join(c.Basepath, c.DiagName) + ".tar.gz"

	tar, err := createNewTar(tarFile)
	if err != nil {
		// handle err
	}
	defer tar.t.Close()
	defer tar.g.Close()

	// tar.Create(TarFile)

	runDockerCmds(tar)
	runSystemCmds(tar)
	runCollectLogs(tar)

	// tar.t.Close()
	// tar.g.Close()

	tar.Finalize(filepath.Join(c.Basepath, c.DiagName+".log"))
	if cfg.UploadUID != "" {
		runUpload(tar.filepath)
	}

	// add an empty line
	fmt.Println()
	return nil
}
