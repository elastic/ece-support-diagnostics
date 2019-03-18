package ecediag

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/elastic/beats/libbeat/logp"
)

// CollectLogs will walk the ElasticFolder path looking for the specific patterns.
//  TODO: needs to be cleaned up and variables should be passed in.
func CollectLogs(tar *Tarball) {
	log := logp.NewLogger("collect-logs")
	log.Info("Collecting ECE log files")
	// TODO: break into concatenated pattern, so the code can be commented.
	elasticLogsPattern := regexp.MustCompile(`\/logs\/|\/zookeeper\/data|ensemble-state-file.json$|stunnel.conf$|replicated.cfg.dynamic$`)
	findPattern(cfg.ElasticFolder, elasticLogsPattern, tar)
}

func findPattern(path string, re *regexp.Regexp, tar *Tarball) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			matched := re.MatchString(filePath)
			if matched && !info.IsDir() {
				tar.AddFile(filePath, info, cfg.ElasticFolder)
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

// func findAll(path string, tar *Tarball) {
// 	err := filepath.Walk(path,
// 		func(filePath string, info os.FileInfo, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			tar.AddFile(filePath, info, cfg.Basepath)
// 			// addFile(filePath, info, basePath, tarball)

// 			return nil
// 		})
// 	if err != nil {
// 		// return nil, err
// 		panic(err)
// 		// log.Println(err)
// 	}
// 	// return files, err
// }
