package ece_support_diagnostic

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func TarLogs() {

	// set up the output file
	file, err := os.Create(filepath.Join(Basepath, "tmp", DiagName) + ".tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	// set up the gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	dockerBase := filepath.Dir(DockerTmpDir)
	systemBase := filepath.Dir(SystemTmpDir)
	findAll(DockerTmpDir, dockerBase, tw)
	findAll(SystemTmpDir, systemBase, tw)

	elasticLogs := filepath.Join(Basepath, "test_data")
	elasticLogsPattern := regexp.MustCompile(`\/logs\/|\/zookeeper\/data`)
	findPattern(elasticLogs, elasticLogsPattern, elasticLogs, tw)
}

func findAll(path string, basePath string, tarball *tar.Writer) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			addFile(filePath, info, basePath, tarball)

			return nil
		})
	if err != nil {
		// return nil, err
		log.Println(err)
	}
	// return files, err
}

func findPattern(path string, re *regexp.Regexp, basePath string, tarball *tar.Writer) {
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			matched := re.MatchString(filePath)
			if matched && !info.IsDir() {
				addFile(filePath, info, basePath, tarball)
			}

			return nil
		})
	if err != nil {
		// return nil, err
		log.Println(err)
	}
	// return files, err
}

func addFile(filePath string, info os.FileInfo, basePath string, tarball *tar.Writer) error {
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = strings.TrimPrefix(filePath, basePath+"/")
	if err := tarball.WriteHeader(header); err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(tarball, file)
	return err
}
