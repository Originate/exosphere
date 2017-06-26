package testHelpers

import (
	"os"
)

func EmptyDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
	f, err := os.Mkdir(dir, os.FileMode(0522))
  if err != nil {
    panic(err)
  }
  f.close()
}
