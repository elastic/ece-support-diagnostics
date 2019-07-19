package uploader

import (
	"runtime"

	"github.com/elastic/eluploader/cmd/eluploader/commands"
)

// RunUpload is used for the Elastic upload service when the `-u {{ upload uui }}` is present
func RunUpload(filePath, uploadUID string) {
	apiURL := "https://upload-staging.elstc.co"
	numWorkers := runtime.NumCPU()
	cmd := &commands.UploadCmd{UploadID: uploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
	cmd.Execute()
}
