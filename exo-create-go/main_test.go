package main_test

import (
	"fmt"
	"html/template"
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

var childOutput string
var cwd string
var appDir string
var appName string
var cmd *exec.Cmd
var err error

func createEmptyApp(cwd, appName string) error {
	appDir = path.Join(cwd, "tmp", appName)
	err := testHelpers.EmptyDir(appDir)
	if err != nil {
		return err
	}
	appConfig := AppConfig{appName, "Empty test application", "1.0.0", "0.22.1"}
	dir, err := os.Executable()
	templatePath := path.Join(dir, "..", "..", "..", "exosphere-shared", "templates", "create-app", "application.yml")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(path.Join(appDir, "application.yml"))
	if err != nil {
		panic(err)
	}
	err = t.Execute(f, appConfig)
	if err != nil {
		panic(err)
	}
	f.Close()
	err = os.Mkdir(path.Join(appDir, ".exosphere"), os.FileMode(0777))
	if err != nil {
		panic(err)
	}
	return nil
}

type AppConfig struct {
	AppName, AppVersion, ExocomVersion, AppDescription string
}

func getCommand(cwd, command string) (*exec.Cmd, error) {
	cmdSegments := strings.Split(path.Join(cwd, "bin", command), " ")
	cmd = exec.Command(cmdSegments[0], cmdSegments[1:]...)
	cmd.Dir = appDir
	if err != nil {
		return cmd, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return cmd, nil
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf("Expected:\n\n%s\n\nto include\n\n%s", haystack, needle)
}

func enterInput(row *gherkin.TableRow) error {
	_, input := row.Cells[0].Value, row.Cells[1].Value
	cmd.Stdin = strings.NewReader(input + "\n")
	return cmd.Run()
}

// nolint gocyclo
func FeatureContext(s *godog.Suite) {

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		return createEmptyApp(cwd, appName)
	})

	s.Step(`^executing "([^"]*)"$`, func(command string) error {
		cmd, err = getCommand(cwd, command)
		if err != nil {
			return err
		}
		return nil
	})

	s.Step(`^starting "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		err := testHelpers.EmptyDir(appDir)
		if err != nil {
			return err
		}
		inputs := make([]string, 4)
		inputs[0] = "my-app-2"
		inputs[1] = "0.0.2"
		inputs[2] = "1.0"
		inputs[3] = "nth"

		cmd, err = getCommand(cwd, command)
		in, err := cmd.StdinPipe()
		cmd.Stdout = os.Stdout
		if err != nil {
			panic(err)
		}
		defer in.Close()
		if err = cmd.Start(); err != nil {
			panic(err)
		}
		for _, input := range inputs {
			fmt.Println(input)
			_, err := in.Write([]byte(input + "\n"))
			if err != nil {
				panic(err)
			}
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
		return nil
	})

	s.Step(`^entering into the wizard:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			err := enterInput(row)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

	s.Step(`^running "([^"]*)" in the terminal$`, func(command string) error {
		appDir = path.Join(cwd, "tmp")
		err := testHelpers.EmptyDir(appDir)
		if err != nil {
			return err
		}
		cmd, err = getCommand(cwd, command)
		if err != nil {
			return err
		}
		outputBytes, err := cmd.Output()
		if err != nil {
			return err
		}
		childOutput = string(outputBytes)
		return nil
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return validateTextContains(childOutput, expectedText)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(text string) error {
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
