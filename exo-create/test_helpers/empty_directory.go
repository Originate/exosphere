package testHelpers

import (
	"os"
)

func EmptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, os.FileMode(0777))
}
