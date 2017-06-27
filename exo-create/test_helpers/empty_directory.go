package testHelpers

import (
	"os"
)

func EmptyDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.Mkdir(dir, os.FileMode(0777))
	if err != nil {
		return err
	}
	return nil
}
