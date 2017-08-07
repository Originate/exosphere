package ostools

import "os"

// DoesFileExist returns true if the file filePath exists, and false otherwise
func DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
