package ecediag

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/elastic/beats/libbeat/logp"
)

func CollectLogs(tar *Tarball) {
	log := logp.NewLogger("collect-logs")
	log.Info("Collecting ECE log files")
	elasticLogsPattern := regexp.MustCompile(`\/logs\/|\/zookeeper\/data`)
	findPattern(ElasticFolder, elasticLogsPattern, tar)
}

func findAll(path string, tar *Tarball) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			tar.AddFile(filePath, info, Basepath)
			// addFile(filePath, info, basePath, tarball)

			return nil
		})
	if err != nil {
		// return nil, err
		panic(err)
		// log.Println(err)
	}
	// return files, err
}

func findPattern(path string, re *regexp.Regexp, tar *Tarball) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			matched := re.MatchString(filePath)
			if matched && !info.IsDir() {
				tar.AddFile(filePath, info, ElasticFolder)
				// addFile(filePath, info, basePath, tarball)
			}

			return nil
		})
	if err != nil {
		// return nil, err
		panic(err)
		// log.Println(err)
	}
	// return files, err
}
