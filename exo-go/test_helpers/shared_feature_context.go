package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/util"
	execplus "github.com/Originate/go-execplus"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var cwd string
var childCmdPlus *execplus.CmdPlus
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

	s.BeforeScenario(func(arg1 interface{}) {
		appDir = ""
	})

	s.AfterSuite(func() {
		if err := os.RemoveAll(appDir); err != nil {
			panic(err)
		}
	})

	s.AfterScenario(func(arg1 interface{}, arg2 error) {
		if childCmdPlus != nil {
			if err := childCmdPlus.Kill(); err != nil {
				panic(err)
			}
			childCmdPlus = nil
		}
		dockerComposeDir := path.Join(appDir, "tmp")
		if util.DoesFileExist(path.Join(dockerComposeDir, "docker-compose.yml")) {
			if err := killTestContainers(dockerComposeDir); err != nil {
				panic(err)
			}
		}
	})

	// Application Setup

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		appDir = path.Join(os.TempDir(), appName)
		return createEmptyApp(appName, cwd)
	})

	s.Step(`^I am in the directory of "([^"]*)" application containing a "([^"]*)" service$`, func(appName, serviceRole string) error {
		appDir = path.Join(os.TempDir(), appName)
		return osutil.CopyRecursively(path.Join(cwd, "..", "example-apps", "test app"), path.Join(os.TempDir(), "test app"))
	})

	// Running / Starting a command

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		var err error
		childOutput, err = util.Run(cwd, command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^running "([^"]*)" in my application directory$`, func(command string) error {
		var err error
		childOutput, err = util.Run(appDir, command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		if err := util.CreateEmptyDirectory(appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
		}
		commandWords, err := util.ParseCommand(command)
		if err != nil {
			return err
		}
		childCmdPlus = execplus.NewCmdPlus(commandWords...)
		childCmdPlus.SetDir(appDir)
		return childCmdPlus.Start()
	})

	s.Step(`^starting "([^"]*)" in my application directory$`, func(command string) error {
		commandWords, err := util.ParseCommand(command)
		if err != nil {
			return err
		}
		childCmdPlus = execplus.NewCmdPlus(commandWords...)
		childCmdPlus.SetDir(appDir)
		return childCmdPlus.Start()
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
		if childCmdPlus != nil {
			return childCmdPlus.WaitForText(text, time.Minute)
		}
		return validateTextContains(childOutput, text)
	})

	s.Step(`^it does not print "([^"]*)" in the terminal$`, func(text string) error {
		if childCmdPlus != nil {
			if err := validateTextContains(childCmdPlus.Output, text); err == nil {
				return fmt.Errorf("Expected the process to not print: %s", text)
			}
			return nil
		}
		return validateTextContains(childOutput, text)
	})

	s.Step(`^I see:$`, func(expectedText *gherkin.DocString) error {
		if childCmdPlus != nil {
			return childCmdPlus.WaitForText(expectedText.Content, time.Second*2)
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
		return childCmdPlus.WaitForText(expectedText, time.Second)
	})

	s.Step(`^I eventually see:$`, func(expectedText *gherkin.DocString) error {
		return childCmdPlus.WaitForText(expectedText.Content, time.Second)
	})

	s.Step(`^waiting until the process ends$`, func() error {
		return childCmdPlus.Wait()
	})

	s.Step(`^it exits with code (\d+)$`, func(expectedExitCode int) error {
		if err := childCmdPlus.Wait(); err != nil {
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

	s.Step(`^my workspace contains the empty directory "([^"]*)"`, func(directory string) error {
		dirPath := path.Join(appDir, directory)
		if !util.IsEmptyDirectory(dirPath) {
			return fmt.Errorf("%s is a not an empty directory", directory)
		}
		return nil
	})

	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my application now contains the file "([^"]*)" with the content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})
}
