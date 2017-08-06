package osplus

import (
	"io"
	"os"
)

// IsEmptyDirectory returns true if dirPath is an empty directory, and false otherwise
func IsEmptyDirectory(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	_, err = f.Readdir(1)
	return err == io.EOF
}
