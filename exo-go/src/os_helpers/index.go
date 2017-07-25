package osHelpers

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/tmrts/boilr/pkg/util/osutil"
)

// DirectoryExists returns true if the directory dirPath is an existing directory,
// and false otherwise
func DirectoryExists(dirPath string) bool {
	return isDirectory(dirPath)
}

// EmptyDir creates an empty dir directory and returns an error if any
func EmptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, os.FileMode(0777))
}

// FileExists returns true if the file filePath exists, and false otherwise
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetSubdirectories returns a slice of subdirectories in the directory dirPath
func GetSubdirectories(dirPath string) (result []string, err error) {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return result, err
	}
	for _, entry := range entries {
		if isDirectory(path.Join(dirPath, entry.Name())) {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

// isDirectory returns true if dirPath is a directory, and false otherwise
func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}

// IsEmpty returns true if dirPath is an empty directory, and false otherwise
func IsEmpty(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	_, err = f.Readdir(1)
	return err == io.EOF
}

// IsEmptyFile returns true if filePath is an empty file, and false otherwise
func IsEmptyFile(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fi.Size() == 0
}

// MoveDir moves srcPath to destPath
func MoveDir(srcPath, destPath string) error {
	err := osutil.CopyRecursively(srcPath, destPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(srcPath)
}

// GetUserHomeDir returns the path to the user's home directory
func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
