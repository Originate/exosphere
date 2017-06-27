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
	"github.com/Originate/exosphere/exo-create-go/test_helpers"
)

var cmd *exec.Cmd
var in io.WriteCloser
var out bytes.Buffer
var appDir string
var err error

type AppConfig struct {
	AppName, AppVersion, ExocomVersion, AppDescription string
}

func run(cwd, command string) error {
	cmdSegments := strings.Split(path.Join(cwd, "bin", command), " ")
	cmd = exec.Command(cmdSegments[0], cmdSegments[1:]...)
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

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf("Expected:\n\n%s\n\nto include\n\n%s", haystack, needle)
}

func enterInput(row *gherkin.TableRow) error {
	_, input := row.Cells[0].Value, row.Cells[1].Value
	_, err := in.Write([]byte(input + "\n"))
	if err != nil {
		fmt.Println("didn't write")
		fmt.Errorf("Error:%s", err)
		return err
	}
	return nil
}

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var cwd, childOutput string
	var err error

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^executing "([^"]*)"$`, func(command string) error {
		return run(cwd, command)
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		err := testHelpers.EmptyDir(appDir)
		if err != nil {
			return err
		}
		return run(cwd, command)
	})

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			err := enterInput(row)
			if err != nil {
				return err
			}
		}
		defer in.Close()
		return nil
	})

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		err := testHelpers.EmptyDir(appDir)
		if err != nil {
			return err
		}
		return run(cwd, command)
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		err = cmd.Wait()
		if err != nil {
			return err
		}
		childOutput = out.String()
		return validateTextContains(childOutput, expectedText)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(text string) error {
		err = cmd.Wait()
		if err != nil {
			return err
		}
		childOutput = out.String()
		return validateTextContains(childOutput, text)
	})

	s.Step(`^my application contains the file "([^"]*)" with the content:$`, func(filePath string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, filePath))
		if err != nil {
			return err
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return err
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my workspace contains the empty directory "([^"]*)"`, func(directory string) error {
		_, err := os.Stat(path.Join(appDir, directory))
		if err == nil {
			return nil
		}
		if os.IsNotExist(err) {
			return err
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
