package ostools

import (
	"io/ioutil"
	"path"
)

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
