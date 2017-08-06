package osplus

// DoesDirectoryExist returns true if the directory dirPath is an existing directory,
// and false otherwise
func DoesDirectoryExist(dirPath string) bool {
	return isDirectory(dirPath)
}
