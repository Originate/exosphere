package testHelpers

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"syscall"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var cwd string
var process *processHelpers.Process
var childOutput string
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

	s.AfterScenario(func(arg1 interface{}, arg2 error) {
		if process != nil {
			if err := process.Kill(); err != nil {
				panic(err)
			}
			process = nil
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
		childOutput, err = processHelpers.Run(cwd, command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^running "([^"]*)" in my application directory$`, func(command string) error {
		var err error
		childOutput, err = processHelpers.Run(appDir, command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		if err := osHelpers.EmptyDir(appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
		}
		commandWords, err := processHelpers.ParseCommand(command)
		if err != nil {
			return err
		}
		process = processHelpers.NewProcess(commandWords...)
		process.SetDir(appDir)
		return process.Start()
	})

	s.Step(`^starting "([^"]*)" in my application directory$`, func(command string) error {
		commandWords, err := processHelpers.ParseCommand(command)
		if err != nil {
			return err
		}
		process = processHelpers.NewProcess(commandWords...)
		process.SetDir(appDir)
		return process.Start()
	})

	// Entering user input

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			if err := enterInput(row); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to enter %s into the wizard", row.Cells[1].Value))
			}
		}
		return nil
	})

	// Verifying output

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(text string) error {
		if process != nil {
			return process.WaitForTextWithTimeout(text, 60000)
		}
		return validateTextContains(childOutput, text)
	})

	s.Step(`^it does not print "([^"]*)" in the terminal$`, func(text string) error {
		if process != nil {
			if err := validateTextContains(process.Output, text); err == nil {
				return fmt.Errorf("Expected the process to not print: %s", text)
			}
			return nil
		}
		return validateTextContains(childOutput, text)
	})

	s.Step(`^I see:$`, func(expectedText *gherkin.DocString) error {
		if process != nil {
			return process.WaitForTextWithTimeout(expectedText.Content, 1500)
		}
		return validateTextContains(childOutput, expectedText.Content)
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
		return process.WaitForTextWithTimeout(expectedText, 1000)
	})

	s.Step(`^I eventually see:$`, func(expectedText *gherkin.DocString) error {
		return process.WaitForTextWithTimeout(expectedText.Content, 1000)
	})

	s.Step(`^waiting until the process ends$`, func() error {
		return process.Wait()
	})

	s.Step(`^it exits with code (\d+)$`, func(expectedExitCode int) error {
		if err := process.Wait(); err != nil {
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
