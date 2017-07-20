package testHelpers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
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
	appDir := path.Join(cwd, "tmp", appName)
	process := processHelpers.NewProcess("exo", "run") // nolint gas
	process.SetDir(appDir)
	err := process.Start()
	if err != nil {
		return err
	}
	return process.WaitForTextWithTimeout("setup complete", 60000)
}

func enterInput(row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := process.WaitForTextWithTimeout(field, 1000); err != nil {
		return err
	}
	_, err := process.StdinPipe.Write([]byte(input + "\n"))
	return err
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}
