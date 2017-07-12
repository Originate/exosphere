package testHelpers

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

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

func enterInput(in io.WriteCloser, out fmt.Stringer, row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := waitForText(out, field, 1000); err != nil {
		return err
	}
	_, err := in.Write([]byte(input + "\n"))
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

func waitForText(stdout fmt.Stringer, text string, duration int) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(time.Duration(duration) * time.Millisecond)
	var output string
	for !strings.Contains(output, text) {
		select {
		case <-ticker.C:
			output = stdout.String()
		case <-timeout:
			return fmt.Errorf("Timed out after %d milliseconds", duration)
		}
	}
	return nil
}
