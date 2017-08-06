package osplus

import "os"

func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}
