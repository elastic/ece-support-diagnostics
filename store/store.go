package store

import "os"

// ContentStore provides the generic storage functions
type ContentStore interface {
	// TODO: remove os.FileInfo, this should be handled internally
	AddFile(filePath string, info os.FileInfo, relPath string) error
	// AddFile(filePath, relativePath string) error

	// REMOVE: the tar file needs the complete count of bytes to properly write the header
	// removing this will force to write a tpmfile, and just pass that in via the AddFile
	AddData(filePath string, b []byte) error
}
