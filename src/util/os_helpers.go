package util

import (
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

// MakeDirectory creates a directory dir if it doesn't already exist, returns an error if any
func MakeDirectory(dir string) error {
	exists, err := DoesDirectoryExist(dir)
	if err != nil {
		return err
	}
	if !exists {
		return os.MkdirAll(dir, os.FileMode(0777))
	}
	return nil
}

// DoesDirectoryExist returns true if the directory dirPath is an existing directory,
// and false otherwise
func DoesDirectoryExist(dirPath string) (bool, error) {
	return isDirectory(dirPath)
}

// DoesFileExist returns true if the file filePath exists, and false otherwise
func DoesFileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
		isDir, err := isDirectory(path.Join(dirPath, entry.Name()))
		if err != nil {
			return result, err
		}
		if isDir {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}

// IsEmptyFile returns true if filePath is an empty file, and false otherwise
func IsEmptyFile(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)
	if err == nil {
		return stat.Size() == 0, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

func isDirectory(dirPath string) (bool, error) {
	stat, err := os.Stat(dirPath)
	if err == nil {
		return stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
