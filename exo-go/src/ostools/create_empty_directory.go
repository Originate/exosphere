package ostools

import "os"

// CreateEmptyDirectory creates an empty dir directory and returns an error if any
func CreateEmptyDirectory(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, os.FileMode(0777))
}
