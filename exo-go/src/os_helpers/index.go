package osHelpers

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/tmrts/boilr/pkg/util/osutil"
)

// IsEmpty returns true if dirPath is an empty directory, and false otherwise
func IsEmpty(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	_, err = f.Readdir(1)
	return err == io.EOF
}

// GetSubdirectories returns a slice of subdirectories in the directory dirPath
func GetSubdirectories(dirPath string) []string {
	subDirectories := []string{}
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if isDirectory(path.Join(dirPath, entry.Name())) {
			subDirectories = append(subDirectories, entry.Name())
		}
	}
	return subDirectories
}

// MoveDir moves srcPath to destPath
func MoveDir(srcPath, destPath string) {
	err := osutil.CopyRecursively(srcPath, destPath)
	if err != nil {
		log.Fatalf("Failed to copy rescursively: %s", err)
	}
	err = os.RemoveAll(srcPath)
	if err != nil {
		log.Fatalf("Failed to remove all: %s", err)
	}
}

// FileExists returns true if the file filePath exists, and false otherwise
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectoryExists returns true if the directory dirPath is an existing directory,
// and false otherwise
func DirectoryExists(dirPath string) bool {
	return isDirectory(dirPath)
}

// isDirectory returns true if dirPath is a directory, and false otherwise
func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}
