package systemLogs

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/config"
	"github.com/elastic/ece-support-diagnostics/store"
)

type fileSystemStore struct {
	store.ContentStore
	cfg *config.Config
}

// File holds file information for a collected log or file
type File struct {
	info     os.FileInfo
	filepath string
}

// Files is an array of the File type, used for holding all the collected logs and files
type Files []File

// Run runs
func Run(t store.ContentStore, config *config.Config) {
	store := fileSystemStore{t, config}
	store.runCollectLogs()
}

// CollectLogs will walk the ElasticFolder path looking for the specific patterns.
//  TODO: needs to be cleaned up and variables should be passed in.
func (t fileSystemStore) runCollectLogs() {
	log := logp.NewLogger("collect-logs")
	log.Infof("Collecting ECE log files")
	// TODO: break into concatenated pattern, so the code can be commented.
	elasticLogsPattern := regexp.MustCompile(`\/logs\/|ensemble-state-file.json$|stunnel.conf$|replicated.cfg.dynamic$`)

	files := Files{}

	t.findPattern(&files, t.cfg.ElasticFolder, elasticLogsPattern)
	// fmt.Printf("%+v\n", files)
	t.addToTar(files)
}

func (t fileSystemStore) findPattern(files *Files, path string, re *regexp.Regexp) {
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

func (t fileSystemStore) addToTar(files Files) {
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
			if t.dateFilter(file, t.cfg.OlderThan) {
				tarRelPath := filepath.Join(t.cfg.DiagName, match[2])
				t.AddFile(file.filepath, file.info, tarRelPath)
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
				if t.dateFilter(file, t.cfg.OlderThan) {
					tarRelPath := filepath.Join(t.cfg.DiagName, "ece", match[2])
					t.AddFile(file.filepath, file.info, tarRelPath)
					l.Infof("Adding log file: %s", match[2])
				}

			} else {
				if file.zeroByteCheck(filepath.Base(file.filepath)) {
					continue
				}
				if t.dateFilter(file, t.cfg.OlderThan) {
					// This should be a catch all. This shouldn't happen.
					l.Warnf("THIS SHOULD NOT HAPPEN, %s", file.filepath)
					tarRelPath := filepath.Join(t.cfg.DiagName, "ece", filepath.Base(file.filepath))
					t.AddFile(file.filepath, file.info, tarRelPath)
					l.Infof("Adding log file: %s", filepath.Base(file.filepath))
				}
			}
		}
	}
}

func (t fileSystemStore) dateFilter(file File, window time.Duration) bool {
	l := logp.NewLogger("files")

	modTime := file.info.ModTime()
	delta := t.cfg.StartTime.Sub(modTime)
	if delta <= window {
		return true
	}
	overLimitString := (delta - window).Truncate(time.Millisecond).String()
	limitString := window.Truncate(time.Millisecond).String()

	fp := strings.TrimLeft(file.filepath, t.cfg.ElasticFolder)

	l.Warnf("Ignoring file: %s threshold, %s too old, %s", limitString, overLimitString, fp)
	return false
}

func (file File) zeroByteCheck(name string) bool {
	l := logp.NewLogger("files")
	if file.info.Size() == int64(0) {
		l.Infof("Skipping 0 byte file: %s", name)
		return true
	}
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
