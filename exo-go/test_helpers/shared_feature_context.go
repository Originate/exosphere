package testHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"syscall"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var cwd string
var cmd *exec.Cmd
var childOutput string
var stdinPipe io.WriteCloser
var stdoutBuffer *bytes.Buffer
var appDir string

// SharedFeatureContext defines the festure context shared between the sub commands
// nolint gocyclo
func SharedFeatureContext(s *godog.Suite) {
	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.AfterSuite(func() {
		if err := os.RemoveAll(appDir); err != nil {
			panic(err)
		}
	})

	// Application Setup

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		appDir = path.Join(os.TempDir(), appName)
		return createEmptyApp(appName, cwd)
	})

	s.Step(`^I am in the directory of "([^"]*)" application containing a "([^"]*)" service$`, func(appName, serviceRole string) error {
		appDir = path.Join(os.TempDir(), appName)
		return osutil.CopyRecursively(path.Join(cwd, "..", "exosphere-shared", "example-apps", "test app"), path.Join(os.TempDir(), "test app"))
	})

	// Running / Starting a command

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		var err error
		childOutput, err = processHelpers.Run(command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^starting "([^"]*)" in my application directory$`, func(command string) error {
		var err error
		cmd, stdinPipe, stdoutBuffer, err = processHelpers.Start(command, appDir)
		return err
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		if err := emptyDir(appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
		}
		var err error
		cmd, stdinPipe, stdoutBuffer, err = processHelpers.Start(command, appDir)
		return err
	})

	// Entering user input

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			if err := enterInput(stdinPipe, stdoutBuffer, row); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to enter %s into the wizard", row.Cells[1].Value))
			}
		}
		return nil
	})

	// Verifying output

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(text string) error {
		if err := validateTextContains(childOutput, text); err != nil {
			return waitForText(stdoutBuffer, text, 5000)
		}
		return nil
	})

	s.Step(`^it does not print "([^"]*)" in the terminal$`, func(text string) error {
		if err := validateTextDoesNotContain(childOutput, text); err != nil {
			return err
		}
		if err := waitForText(stdoutBuffer, text, 1500); err != nil {
			return nil
		}
		return fmt.Errorf("expect the process to not print: %s", text)
	})

	s.Step(`^I see:$`, func(docString *gherkin.DocString) error {
		return validateTextContains(childOutput, docString.Content)
	})

	s.Step(`^the output matches "([^"]*)"$`, func(text string) error {
		matched, err := regexp.Match(text, []byte(childOutput))
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("output does not match")
		}
		return nil
	})

	s.Step(`^I eventually see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return waitForText(stdoutBuffer, expectedText, 1000)
	})

	s.Step(`^I eventually see:$`, func(expectedText *gherkin.DocString) error {
		return waitForText(stdoutBuffer, expectedText.Content, 1000)
	})

	s.Step(`^waiting until the process ends$`, func() error {
		return cmd.Wait()
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
