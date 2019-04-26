package store

import "os"

type ContentStore interface {
	// AddFile(filePath, relativePath string) error
	AddFile(filePath string, info os.FileInfo, relPath string) error
	AddData(filePath string, b []byte) error
}
