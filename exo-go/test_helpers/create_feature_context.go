package testHelpers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
)

func enterInput(in io.WriteCloser, out bytes.Buffer, row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := waitForText(out, field, 1000); err != nil {
		return err
	}
	_, err := in.Write([]byte(input + "\n"))
	return err
}

// CreateFeatureContext defines the festure context for features/create.feature
// nolint gocyclo
func CreateFeatureContext(s *godog.Suite) {
	var cwd string
	var stdinPipe io.WriteCloser
	var stdoutBuffer bytes.Buffer
	var appDir string

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		if err := emptyDir(appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
		}
		var err error
		stdinPipe, stdoutBuffer, err = start(command, appDir)
		return err
	})

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			if err := enterInput(stdinPipe, stdoutBuffer, row); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to enter %s into the wizard", row.Cells[1].Value))
			}
		}
		return nil
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return waitForText(stdoutBuffer, expectedText, 1000)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(expectedText string) error {
		return waitForText(stdoutBuffer, expectedText, 1000)
	})

	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my workspace contains the empty directory "([^"]*)"`, func(directory string) error {
		f, err := os.Stat(path.Join(appDir, directory))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to get information for the entry %s", directory))
		}
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", directory)
		} else if !f.IsDir() {
			return fmt.Errorf("%s is a not a directory", directory)
		}
		return nil
	})
}
