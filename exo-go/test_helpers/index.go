package testHelpers

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

func checkoutApp(cwd, appName string) error {
	src := path.Join(cwd, "..", "exosphere-shared", "example-apps", appName)
	dest := path.Join(cwd, "tmp", appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return CopyDir(src, dest)
}

func setupApp(cwd, appName string) error {
	cmdPath := path.Join(cwd, "..", "exo-setup", "bin", "exo-setup")
	cmd := exec.Command(cmdPath) // nolint gas
	cmd.Dir = path.Join(cwd, "tmp", appName)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running setup\nOutput:\n%s\nError:%s", string(outputBytes), err)
	}
	return nil
}

func enterInput(row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := process.WaitForTextWithTimeout(field, 1000); err != nil {
		return err
	}
	_, err := process.StdinPipe.Write([]byte(input + "\n"))
	return err
}

func emptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, os.FileMode(0777))
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}
