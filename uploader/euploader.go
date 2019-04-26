package uploader

import (
	"runtime"

	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/eluploader/cmd/eluploader/commands"
)

// runUpload is used for the Elastic upload service when the `-u {{ upload uui }}` is present
func RunUpload(filePath string, cfg *config.Config) {
	apiURL := "https://upload-staging.elstc.co"
	numWorkers := runtime.NumCPU()
	cmd := &commands.UploadCmd{UploadID: cfg.UploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
	cmd.Execute()
}
