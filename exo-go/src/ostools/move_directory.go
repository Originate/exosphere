package ostools

import (
	"os"

	"github.com/tmrts/boilr/pkg/util/osutil"
)

// MoveDirectory moves srcPath to destPath
func MoveDirectory(srcPath, destPath string) error {
	err := osutil.CopyRecursively(srcPath, destPath)
	if err != nil {
		return err
	}
	return os.RemoveAll(srcPath)
}
