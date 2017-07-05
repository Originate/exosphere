package osHelpers

import (
	"io"
	"io/ioutil"
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
	if err == io.EOF {
		return true
	}
	return false
}

// IsValidTemplateDir returns true if the directory templateDir is a valid
// boilr template directory
func IsValidTemplateDir(templateDir string) bool {
	return FileExists(path.Join(templateDir, "project.json")) && DirectoryExists(path.Join(templateDir, "template"))
}

// GetSubdirectories returns a slice of subdirectories in the directory dirPath
func GetSubdirectories(dirPath string) []string {
	subDirectories := []string{}
	entries, _ := ioutil.ReadDir(dirPath)
	for _, entry := range entries {
		if isDirectory(path.Join(dirPath, entry.Name())) {
			subDirectories = append(subDirectories, entry.Name())
		}
	}
	return subDirectories
}

// MoveDir moves srcPath to destPath
func MoveDir(srcPath, destPath string) {
	osutil.CopyRecursively(srcPath, destPath)
	os.RemoveAll(srcPath)
}

// FileExists returns true if the file filePath exists, and false otherwise
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DirectoryExists returns true if the directory dirPath is an existing directory,
// and false otherwise
func DirectoryExists(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}

// isDirectory returns true if dirPath is a directory, and false otherwise
func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}
