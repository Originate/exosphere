package osHelpers

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/tmrts/boilr/pkg/util/osutil"
)

func IsEmpty(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	_, err = f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

func IsValidTemplateDir(templateDir string) bool {
	return FileExists(path.Join(templateDir, "project.json")) && DirectoryExists(path.Join(templateDir, "template"))
}

func GetSubdirectories(dirPath string) []string {
	subDirectories := []string{}
	entries, _ := ioutil.ReadDir(dirPath)
	for _, entry := range entries {
		if isDirectory(path.Join(dirPath, entry.Name())) {
			subDirectories = append(subDirectories, entry.Name())
		}
	}
	return subDirectories
}

func MoveDir(srcPath, destPath string) {
	osutil.CopyRecursively(srcPath, destPath)
	os.RemoveAll(srcPath)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func DirectoryExists(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}

func isDirectory(dirPath string) bool {
	fi, err := os.Stat(dirPath)
	return err == nil && fi.IsDir()
}
