package appSetup

import (
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/test_helpers"
)

// CheckoutApp copies the example app appName to cwd
func CheckoutApp(cwd, appName string) error {
	src := path.Join("..", "..", "..", "exosphere-shared", "example-apps", appName)
	dest := path.Join(cwd, appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return testHelpers.CopyDir(src, dest)
}
