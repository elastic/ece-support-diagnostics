package uploader

import (
	"runtime"

	"github.com/elastic/eluploader/cmd/eluploader/commands"
)

// RunUpload is used for the Elastic upload service when the `-U {{ upload uuid }}` is present
func RunUpload(filePath, uploadUID string) {
	apiURL := "https://upload.elastic.co"
	numWorkers := runtime.NumCPU()
	cmd := &commands.UploadCmd{UploadID: uploadUID, Filepath: filePath, ApiURL: apiURL, NumWorkers: numWorkers}
	cmd.Execute()
}
