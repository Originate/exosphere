package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/src"
	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
	"github.com/pkg/errors"
)

var childCmdPlus *execplus.CmdPlus
var childOutput string
var appDir string
var templateDir string

func waitWithTimeout(cmdPlus *execplus.CmdPlus, duration time.Duration) error {
	done := make(chan error)
	go func() { done <- cmdPlus.Wait() }()
	select {
	case err := <-done:
		return err
	case <-time.After(duration):
		return fmt.Errorf("Timed out after %v, command did not exit. Full output:\n%s", duration, cmdPlus.GetOutput())
	}
}

// SharedFeatureContext defines the festure context shared between the sub commands
// nolint: gocyclo
func SharedFeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(arg1 interface{}) {
		var err error
		appDir, err = ioutil.TempDir("", "")
		if err != nil {
			panic(err)
		}
		templateDir, err = ioutil.TempDir("", "")
		if err != nil {
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
		if err := cleanApp(appDir); err != nil {
			panic(err)
		}
	})

	// Application Setup

	s.Step(`^I am in the root directory of a non-exosphere application$`, func() error {
		return nil
	})

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		var err error
		appDir, err = createEmptyApp(appName)
		return err
	})

	s.Step(`^I am in the root directory of the "([^"]*)" example application$`, func(name string) error {
		return CheckoutApp(appDir, name)
	})

	s.Step(`^it doesn\'t run any tests$`, func() error {
		expectedText := "is not an exosphere application"
		if childCmdPlus != nil {
			return childCmdPlus.WaitForText(expectedText, time.Minute)
		}
		return validateTextContains(childOutput, expectedText)
	})

	s.Step(`^I am in the directory of "([^"]*)" application containing a "([^"]*)" service$`, func(appName, serviceRole string) error {
		return CheckoutApp(appDir, "test-app")
	})

	s.Step(`^my application has the templates:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			templateName, gitURL := row.Cells[0].Value, row.Cells[1].Value
			command := []string{"exo", "template", "add", templateName, gitURL}
			if len(row.Cells) == 3 {
				gitTag := row.Cells[2].Value
				command = append(command, gitTag)
			}
			if _, err := util.Run(appDir, command...); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to creates the template %s:%s\n", appDir, err))
			}
		}
		return nil
	})

	s.Step(`^my (?:application|workspace) contains the empty directory "([^"]*)"`, func(directory string) error {
		dirPath := path.Join(appDir, directory)
		isDir, err := util.DoesDirectoryExist(dirPath)
		if err != nil {
			return err
		}
		if !isDir {
			return fmt.Errorf("%s is a not a directory", dirPath)
		}
		fileInfos, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return err
		}
		if len(fileInfos) != 0 {
			fileNames := []string{}
			for _, fileInfo := range fileInfos {
				fileNames = append(fileNames, fileInfo.Name())
			}
			return fmt.Errorf("%s is a not a an empty directory. Contains: %s", dirPath, strings.Join(fileNames, ", "))
		}
		return nil
	})

	// Running / Starting a command

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		var err error
		childOutput, err = util.Run(appDir, command)
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

	s.Step(`^running "([^"]*)" in the "([^"]*)" directory$`, func(command, dirName string) error {
		var err error
		childOutput, err = util.Run(path.Join(appDir, dirName), command)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command errored with output: %s", childOutput))
		}
		return nil
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
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

	s.Step(`^starting "([^"]*)" in the "([^"]*)" directory$`, func(command, dirName string) error {
		commandWords, err := util.ParseCommand(command)
		if err != nil {
			return err
		}
		childCmdPlus = execplus.NewCmdPlus(commandWords...)
		childCmdPlus.SetDir(path.Join(appDir, dirName))
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

	s.Step(`^it prints the current version in the terminal$`, func() error {
		text := fmt.Sprintf("Exosphere v%s", src.Version)
		if childCmdPlus != nil {
			return childCmdPlus.WaitForText(text, time.Minute)
		}
		return validateTextContains(childOutput, text)
	})

	s.Step(`^it does not print "([^"]*)" in the terminal$`, func(text string) error {
		if childCmdPlus != nil {
			if err := validateTextContains(childCmdPlus.GetOutput(), text); err == nil {
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

	s.Step(`^I eventually see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return childCmdPlus.WaitForText(expectedText, time.Second)
	})

	s.Step(`^I eventually see:$`, func(expectedText *gherkin.DocString) error {
		return childCmdPlus.WaitForText(expectedText.Content, time.Minute)
	})

	s.Step(`^(it exits|waiting until the process ends)$`, func() error {
		return waitWithTimeout(childCmdPlus, time.Minute)
	})

	s.Step(`^I stop all running processes$`, func() error {
		if childCmdPlus != nil {
			err := childCmdPlus.Cmd.Process.Signal(os.Interrupt)
			if err != nil {
				return err
			}
			err = waitWithTimeout(childCmdPlus, time.Minute)
			if err != nil {
				fmt.Println("Command did not exit after 1m (TODO: fix this)")
				return childCmdPlus.Kill()
			}
		}
		return nil
	})

	s.Step(`^it exits with code (\d+)$`, func(expectedExitCode int) error {
		actualExitCode := 0
		if err := waitWithTimeout(childCmdPlus, time.Minute); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					actualExitCode = status.ExitStatus()
				} else {
					return fmt.Errorf("Unable to parse Status object: %v", err)
				}
			} else {
				return fmt.Errorf("cmd.Wait: %v", err)
			}
		}
		if actualExitCode != expectedExitCode {
			return fmt.Errorf("Exited with code %d instead of %d. Output:\n%s", actualExitCode, expectedExitCode, childCmdPlus.GetOutput())
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
