package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/ece-support-diagnostics/helpers"
)

// Tarball provides a wrapper for the tar/gz writers, and a mutex lock to call for thread safety
type Tarball struct {
	f    *os.File
	tar  *tar.Writer
	gzip *gzip.Writer
	m    sync.Mutex
}

// Create starts a new tar/gz file to write data into
func Create(filePath string) (*Tarball, error) {
	return createNewTar(filePath)
}

// Filepath provides the full filepath to the tar file
func (tw *Tarball) Filepath() string {
	fp, _ := filepath.Abs(tw.f.Name())
	return fp
}

// Filename provides the filename of the tar file
func (tw *Tarball) Filename() string {
	return tw.f.Name()
}

// Close will finalize the tar / gzip file.
func (tw *Tarball) Close() {
	tw.tar.Close()
	tw.gzip.Close()
	tw.f.Close()
}

func createNewTar(tarballFilePath string) (*Tarball, error) {
	t := new(Tarball)

	file, err := os.Create(tarballFilePath)
	if err != nil {
		return t, fmt.Errorf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error())
	}
	t.f = file
	t.gzip = gzip.NewWriter(file)
	t.tar = tar.NewWriter(t.gzip)

	return t, nil
}

// Finalize adds the logfile to the tar, and closes the tar.
func (tw *Tarball) Finalize(logfilePath, tarRelPath string) {
	l := logp.NewLogger("TarFile")
	l.Infof("Adding log file: %s", logfilePath)

	msgClosingTar := fmt.Sprintf(" the tar: %s", tw.Filepath())
	l.Infof("Finalizing %s", msgClosingTar)
	fmt.Println("[ ] Finalizing" + msgClosingTar)

	fileInfo, err := os.Stat(logfilePath)
	helpers.PanicError(err)

	tw.AddFile(logfilePath, fileInfo, tarRelPath)
	tw.Close()

	helpers.ClearStdoutLine()
	fmt.Println("[✔] Finished" + msgClosingTar)
}

// AddData is for adding byte data directly to the tar file
// Need to figure out how to consume the bytes as streaming io.Writer
func (tw *Tarball) AddData(filePath string, b []byte) error {
	tw.m.Lock()
	defer tw.m.Unlock()

	// Make sure the path does not start with a slash
	filePath = strings.TrimLeft(filePath, "/")

	header := &tar.Header{
		Name:    filePath,
		Size:    int64(len(b)),
		Mode:    int64(0644),
		ModTime: time.Now(),
	}
	err := tw.tar.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("Could not write header for file '%s', got error '%s'", filePath, err.Error())
	}
	tw.tar.Write(b)
	return err
}

// AddFile reads a file and adds it to the tar file. The basePath is removed from the filepath for
//  the path preserved in the tar file.
// func (tw *Tarball) AddFile(filePath string, info os.FileInfo, basePath string) error {
func (tw *Tarball) AddFile(filePath string, info os.FileInfo, relPath string) error {
	tw.m.Lock()
	defer tw.m.Unlock()

	// fmt.Println(filePath)
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// archiveFile := strings.TrimLeft(strings.TrimPrefix(filePath, strings.TrimRight(basePath, "/")), "/")
	// archiveFilePath := filepath.Join(cfg.DiagnosticFilename(), archiveFile)
	// header.Name = archiveFilePath
	header.Name = relPath
	// fmt.Println(header.Name)

	err = tw.tar.WriteHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = io.Copy(tw.tar, file)
	return err
}
