package util

import (
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/tmrts/boilr/pkg/util/osutil"
)

// CreateEmptyDirectory creates an empty dir directory and returns an error if any
func CreateEmptyDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, os.FileMode(0777))
}

// DoesDirectoryExist returns true if the directory dirPath is an existing directory,
// and false otherwise
func DoesDirectoryExist(dirPath string) bool {
	return isDirectory(dirPath)
}

// DoesFileExist returns true if the file filePath exists, and false otherwise
func DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// GetHomeDirectory returns the path to the user's home directory
func GetHomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
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

// IsEmptyDirectory returns true if dirPath is an empty directory, and false otherwise
func IsEmptyDirectory(dirPath string) bool {
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

// MoveDirectory moves srcPath to destPath
func MoveDirectory(srcPath, destPath string) error {
	err := osutil.CopyRecursively(srcPath, destPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(srcPath)
}

// Helpers

func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}
