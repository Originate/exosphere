package testHelpers

import (
	"os"
)

func EmptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.Mkdir(dir, os.FileMode(0777)); err != nil {
		return err
	}
	return nil
}
