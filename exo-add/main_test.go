package main_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-add/test_helpers"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var cmd *exec.Cmd
var in io.WriteCloser
var out bytes.Buffer
var appDir string
var err error

// Runs the given command in the given directory
func run(command []string, dir string) error {
	cmd = exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	in, err = cmd.StdinPipe()
	out.Reset()
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

func createEmptyApp(appName, cwd string) error {
	appDir = path.Join(cwd, "tmp")
	if err := testHelpers.EmptyDir(appDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
	}
	command := path.Join("..", "..", "exo-create", "bin", "exo-create")
	if err := run([]string{command}, path.Join(cwd, "tmp")); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s application directory", appDir))
	}
	fields := []string{"AppName", "AppDescription", "AppVersion", "ExocomVersion"}
	inputs := []string{appName, "Empty test application", "1.0.0", "0.22.1"}
	for i, field := range fields {
		if err = testHelpers.WaitForText(&out, field, 1000); err != nil {
			return err
		}
		if _, err := in.Write([]byte(inputs[i] + "\n")); err != nil {
			return err
		}
	}
	return nil
}

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var cwd, tmpDir, appDir string

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		tmpDir = path.Join(cwd, "tmp")
		if err := testHelpers.EmptyDir(tmpDir); err != nil {
			panic(err)
		}
	})

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		appDir = path.Join(tmpDir, appName)
		return createEmptyApp(appName, cwd)
	})

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		splitCommand := strings.Split(command, " ")
		parsedCommand := append([]string{path.Join("..", "..", splitCommand[0], "bin", splitCommand[0])}, splitCommand[1:]...)
		return run(parsedCommand, tmpDir)
	})

	s.Step(`^starting "([^"]*)" in this application's directory$`, func(command string) error {
		return run([]string{path.Join("..", "..", "..", command, "bin", command)}, appDir)
	})

	s.Step(`^my application contains the template folder "([^"]*)" with the files:$`, func(dirPath string, table *gherkin.DataTable) error {
		if err := os.MkdirAll(path.Join(appDir, dirPath), 0777); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create %s", dirPath))
		}
		for _, row := range table.Rows[1:] {
			file, content := row.Cells[0].Value, row.Cells[1].Value
			match, err := regexp.MatchString("/", file)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to parse %s", file))
			}
			if match {
				if err := os.MkdirAll(path.Join(appDir, dirPath, filepath.Dir(file)), 0777); err != nil {
					return errors.Wrap(err, fmt.Sprintf("Failed to create the necessary directories for %s", file))
				}
			}
			if err := ioutil.WriteFile(path.Join(appDir, dirPath, file), []byte(content), 0777); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to create %s", file))
			}
		}
		return nil
	})

	s.Step(`^I am in the directory of "([^"]*)" application containing a "([^"]*)" service$`, func(appName, serviceRole string) error {
		appDir = path.Join(tmpDir, appName)
		return osutil.CopyRecursively(path.Join(cwd, "..", "exosphere-shared", "example-apps", "test app"), path.Join(tmpDir, "test app"))
	})

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			if err := enterInput(row); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to enter %s into the wizard", row.Cells[1].Value))
			}
		}
		return nil
	})

	s.Step(`^waiting until the process ends$`, func() error {
		return testHelpers.WaitForText(&out, "done", 1000)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(expectedText string) error {
		return testHelpers.WaitForText(&out, expectedText, 1000)
	})

	s.Step(`^I see:$`, func(expectedText *gherkin.DocString) error {
		return testHelpers.WaitForText(&out, expectedText.Content, 1000)
	})

	s.Step(`^my application contains the file "([^"]*)" with the content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return testHelpers.ValidateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my application contains the file "([^"]*)" containing the text:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return testHelpers.ValidateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my application contains the empty directory "([^"]*)"`, func(directory string) error {
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

	s.Step(`^I see the error "([^"]*)"$`, func(expectedText string) error {
		// fmt.Print(expectedText)
		return testHelpers.WaitForText(&out, expectedText, 1000)
	})

	s.Step(`^it exits with code (\d+)$`, func(expectedExitCode int) error {
		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					if status.ExitStatus() != expectedExitCode {
						return fmt.Errorf("Exited with code %d instead of %d", status.ExitStatus(), expectedExitCode)
					}
					return nil
				}
			} else {
				return fmt.Errorf("cmd.Wait: %v", err)
			}
		}
		return fmt.Errorf("Expected to exit with code: %d", expectedExitCode)
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
