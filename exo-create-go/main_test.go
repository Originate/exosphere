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
)

func createEmptyApp(cwd, appName string) error {
	appDir = path.Join(cwd, "tmp", appName)
	err := testHelpers.EmptyDir(appDir)
	if err != nil {
		return err
	}
	appConfig := AppConfig{appName, "Empty test application", "1.0.0", "0.22.1"}
	templatePath := path.Join("..", "exosphere-shared", "templates", "create-app", "application.yml")
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
	f.close()
	f, err := os.Mkdir(path.Join(appDir, ".exosphere"), os.FileMode(0522))
	if err != nil {
		panic(err)
	}
	f.close()
	return nil
}

type AppConfig struct {
	AppName, AppVersion, ExocomVersion, AppDescription string
}

func run(cwd, command, appName string) (Cmd, error) {
	cmdPath := path.Join(cwd, "bin", command)
	cmd := exec.Command(cmdPath)
	cmd.Dir = path.Join(cwd, "tmp", appName)
	err := cmd.Run()
	if err != nil {
		return cmd, fmt.Errorf("Error running %s\nOutput:\n%s\nError:%s", command, string(outputBytes), err)
	}
	return cmd, nil
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf("Expected:\n\n%s\n\nto include\n\n%s", haystack, needle)
}

func enterInput(row []string) error {
	_, input := row[0], row[1]
	cmd.Stdin = strings.NewReader("%s\n", input)
	return cmd.Run()
}

// nolint gocyclo
func FeatureContext(s *godog.Suite) {
	var childOutput string
	var cwd string
	var appDir string
	var cmd Cmd

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^I am in the root directory of an empty application called "([^"]*)"$`, func(appName string) error {
		createEmptyApp(cwd, appName)
	})

	s.Step(`^executing "([^"]*)"$`, func(command string) error {
		cmd, err = run(cwd, command, appName)
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
		cmd, err = run(cwd, command, appName)
		if err != nil {
			return err
		}
		childOutput = string(cmd.CombinedOutput())
		return nil
	})

	s.Step(`^entering into the wizard:$`, func(table [][]string, command string) error {
		for row := range table.rows() {
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
		cmd, err = run(cwd, command, appName)
		if err != nil {
			return err
		}
		childOutput = string(cmd.CombinedOutput())
		return nil
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return validateTextContains(childOutput, expectedText)
	})

	s.Step(`^it prints "([^"]*)" in the terminal$`, func(text string) error {
		return validateTextContains(childOutput, text)
	})

	s.Step(`^my application contains the file "([^"]*)" with the content:$`, func(filePath, expectedContent string) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, filePath))
		return strings.Contains(strings.TrimSpaces(string(bytes)), strings.TrimSpaces(expectedContent))
	})

	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName, expectedContent string) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, filename))
		return strings.Contains(strings.TrimSpaces(string(bytes)), strings.TrimSpaces(expectedContent))
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
