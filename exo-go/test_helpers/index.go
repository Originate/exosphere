package testHelpers

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

func enterInput(in io.WriteCloser, out io.ReadCloser, row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := process.WaitForText(field, 1000); err != nil {
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
