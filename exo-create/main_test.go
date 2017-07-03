package main_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-create/test_helpers"
	"github.com/pkg/errors"
)

var cmd *exec.Cmd
var in io.WriteCloser
var out bytes.Buffer
var appDir string
var err error

// Runs the given command in the given directory
func run(command []string) error {
	cmd = exec.Command(command[0], command[1:]...)
	cmd.Dir = appDir
	in, err = cmd.StdinPipe()
	cmd.Stdout = &out
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return nil
}

func enterInput(row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err = testHelpers.WaitForText(&out, field, 1000); err != nil {
		return err
	}
	_, err := in.Write([]byte(input + "\n"))
	return err
}

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var cwd string

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		if err := testHelpers.EmptyDir(appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
		}
		return run(testHelpers.ReformatCommand(command, cwd))
	})

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			if err := enterInput(row); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to enter %s into the wizard", row.Cells[1].Value))
			}
		}
		return nil
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return testHelpers.WaitForText(&out, expectedText, 1000)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(expectedText string) error {
		return testHelpers.WaitForText(&out, expectedText, 1000)
	})

	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return testHelpers.ValidateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
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

func TestMain(m *testing.M) {
	var paths []string
	var format string
	if len(os.Args) == 3 && os.Args[1] == "--" {
		format = "pretty"
		paths = append(paths, os.Args[2])
	} else {
		format = "progress"
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:        format,
		NoColors:      false,
		StopOnFailure: true,
		Paths:         paths,
	})
	os.Exit(status)
}
