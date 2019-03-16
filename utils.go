package ecediag

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elastic/beats/libbeat/logp"
)

func setupFolders(folders []string) {
	for _, folder := range folders {
		f, _ := filepath.Abs(folder)
		flog := logp.NewLogger("folders")
		flog.Info("Temp folder setup: ", f)
		os.MkdirAll(f, os.ModePerm)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func writeFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, 0644)
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func clearStdoutLine() {
	fmt.Printf("\033[F") // back to previous line
	fmt.Printf("\033[K") // clear line
}
