package appSetup

import (
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
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

func getCommands(dockerfilePath string) ([]string, error) {
	var commands []string
	if osHelpers.FileExists(dockerfilePath) {
		f, err := os.Open(dockerfilePath)
		if err != nil {
			return commands, err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "RUN") {
				commands = append(commands, strings.Split(line, "RUN ")[1])
			}
		}
		if err := f.Close(); err != nil {
			return commands, err
		}
	}
	return commands, nil
}
