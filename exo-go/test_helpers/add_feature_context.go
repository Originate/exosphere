package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
)

func createEmptyApp(appName, cwd string) error {
	appDir = path.Join(os.TempDir(), appName)
	if err := emptyDir(appDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create an empty %s directory", appDir))
	}
	process = processHelpers.NewProcess("exo", "create")
	process.SetDir(os.TempDir())
	if err := process.Start(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create '%s' application", appDir))
	}
	fields := []string{"AppName", "AppDescription", "AppVersion", "ExocomVersion"}
	inputs := []string{appName, "Empty test application", "1.0.0", "0.22.1"}
	for i, field := range fields {
		if err := process.WaitForTextWithTimeout(field, 5000); err != nil {
			return err
		}
		if _, err := process.StdinPipe.Write([]byte(inputs[i] + "\n")); err != nil {
			return err
		}
	}
	return nil
}

// nolint gocyclo
func AddFeatureContext(s *godog.Suite) {
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

	s.Step(`^my application now contains the file "([^"]*)" with the content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my application now contains the file "([^"]*)" containing the text:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
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
}
