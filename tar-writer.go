package ecediag

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
)

// Tarball provides a wrapper for the tar/gz writers, and a mutex lock to call for thread safety
type Tarball struct {
	filepath string
	t        *tar.Writer
	g        *gzip.Writer
	m        sync.Mutex
}

// Create starts a new tar/gz file to write data into
func (tw *Tarball) Create(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	// // set up the output file
	// defer file.Close()
	tw.filepath = filePath

	// set up the gzip writer
	gw := gzip.NewWriter(file)
	tw.g = gw
	// defer gw.Close()

	t := tar.NewWriter(gw)
	tw.t = t
	return err
}

// Finalize adds the logfile to the tar, and closes the tar.
func (tw *Tarball) Finalize(logfilePath string) {

	// TODO: This needs to be improved. I would like to just call AddFile.
	//  need to make Addfile take a struct that has the stat,name,tar filepath,etc
	l := logp.NewLogger("TarFile")
	l.Infof("Adding log file: %s", logfilePath)

	msgClosingTar := fmt.Sprintf(" the tar: %s", tw.filepath)
	l.Infof("Finalizing %s", msgClosingTar)
	fmt.Println("[ ] Finalizing" + msgClosingTar)

	fileInfo, err := os.Stat(logfilePath)
	panicError(err)

	logTarPath := filepath.Join(DiagName, "diagnostic.log")

	tw.m.Lock()
	defer tw.m.Unlock()

	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	panicError(err)
	header.Name = logTarPath

	err = tw.t.WriteHeader(header)
	panicError(err)

	file, err := os.Open(logfilePath)
	panicError(err)
	defer file.Close()

	_, err = io.Copy(tw.t, file)
	panicError(err)

	tw.t.Close()
	tw.g.Close()

	clearStdoutLine()
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
	err := tw.t.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("Could not write header for file '%s', got error '%s'", filePath, err.Error())
	}
	tw.t.Write(b)
	return err
}

// AddFile reads a file and adds it to the tar file. The basePath is removed from the filepath for
//  the path preserved in the tar file.
func (tw *Tarball) AddFile(filePath string, info os.FileInfo, basePath string) error {
	tw.m.Lock()
	defer tw.m.Unlock()

	// fmt.Println(filePath)
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	archiveFile := strings.TrimLeft(strings.TrimPrefix(filePath, strings.TrimRight(basePath, "/")), "/")
	archiveFilePath := filepath.Join(DiagName, archiveFile)
	header.Name = archiveFilePath
	// fmt.Println(header.Name)

	err = tw.t.WriteHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = io.Copy(tw.t, file)
	return err
}

// https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185
// type chanWriter struct {
// 	ch chan byte
// }
//
// func newChanWriter() *chanWriter {
// 	return &chanWriter{make(chan byte, 1024)}
// }
//
// func (w *chanWriter) Chan() <-chan byte {
// 	return w.ch
// }
//
// func (w *chanWriter) Write(p []byte) (int, error) {
// 	n := 0
// 	for _, b := range p {
// 		w.ch <- b
// 		n++
// 	}
// 	return n, nil
// }
//
// func (w *chanWriter) Close() error {
// 	close(w.ch)
// 	return nil
// }
