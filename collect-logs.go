package ecediag

import (
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

// File holds file information for a collected log or file
type File struct {
	info     os.FileInfo
	filepath string
}

// Files is an array of the File type, used for holding all the collected logs and files
type Files []File

// CollectLogs will walk the ElasticFolder path looking for the specific patterns.
//  TODO: needs to be cleaned up and variables should be passed in.
func runCollectLogs(tar *Tarball) {
	log := logp.NewLogger("collect-logs")
	log.Infof("Collecting ECE log files")
	// TODO: break into concatenated pattern, so the code can be commented.
	elasticLogsPattern := regexp.MustCompile(`\/logs\/|ensemble-state-file.json$|stunnel.conf$|replicated.cfg.dynamic$`)

	files := Files{}

	files.findPattern(cfg.ElasticFolder, elasticLogsPattern)
	// fmt.Printf("%+v\n", files)
	files.addToTar(tar)
}

func (files *Files) findPattern(path string, re *regexp.Regexp) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			matched := re.MatchString(filePath)
			if matched && !info.IsDir() {
				f := File{info: info, filepath: filePath}
				*files = append(*files, f)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
}

func (files Files) addToTar(tar *Tarball) {
	l := logp.NewLogger("files")
	clusterDiskPathRegex := regexp.MustCompile(`(.*\/services\/allocator\/containers\/)((?:elasticsearch|kibana).*)`)
	eceDiskPathRegex := regexp.MustCompile(`(.*(?:\/services\/|bootstrap-logs\/))(.*)$`)

	for _, file := range files {
		match := clusterDiskPathRegex.FindStringSubmatch(file.filepath)

		if len(match) == 3 {
			// strip off *elastic/172.16.0.25/services/allocator/containers/
			//  this should match elasticsearch and kibana clusters
			if file.zeroByteCheck(match[2]) {
				continue
			}
			if file.dateFilter(cfg.OlderThan) {
				tarRelPath := filepath.Join(cfg.DiagName, match[2])
				tar.AddFile(file.filepath, file.info, tarRelPath)
				l.Infof("Adding log file: %s", match[2])
			}

		} else {
			// strip off *elastic/172.16.0.25/services/ or *elastic/logs/bootstrap-logs/
			//  everything after will be the path stored in the tar
			match := eceDiskPathRegex.FindStringSubmatch(file.filepath)
			if len(match) == 3 {
				if file.zeroByteCheck(match[2]) {
					continue
				}
				if file.dateFilter(cfg.OlderThan) {
					tarRelPath := filepath.Join(cfg.DiagName, "ece", match[2])
					tar.AddFile(file.filepath, file.info, tarRelPath)
					l.Infof("Adding log file: %s", match[2])
				}

			} else {
				if file.zeroByteCheck(filepath.Base(file.filepath)) {
					continue
				}
				if file.dateFilter(cfg.OlderThan) {
					// This should be a catch all. This shouldn't happen.
					l.Warnf("THIS SHOULD NOT HAPPEN, %s", file.filepath)
					tarRelPath := filepath.Join(cfg.DiagName, "ece", filepath.Base(file.filepath))
					tar.AddFile(file.filepath, file.info, tarRelPath)
					l.Infof("Adding log file: %s", filepath.Base(file.filepath))
				}
			}
		}
	}
}

func (file File) zeroByteCheck(name string) bool {
	l := logp.NewLogger("files")
	if file.info.Size() == int64(0) {
		l.Infof("Skipping 0 byte file: %s", name)
		return true
	}
	return false
}

func (file File) dateFilter(window time.Duration) bool {
	l := logp.NewLogger("files")

	modTime := file.info.ModTime()
	delta := cfg.StartTime.Sub(modTime)
	if delta <= window {
		return true
	}
	l.Warnf("Ignoring file: %s, %s too old", filepath.Base(file.filepath), delta-window)
	return false
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
