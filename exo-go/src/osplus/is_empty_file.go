package osplus

import "os"

// IsEmptyFile returns true if filePath is an empty file, and false otherwise
func IsEmptyFile(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fi.Size() == 0
}
