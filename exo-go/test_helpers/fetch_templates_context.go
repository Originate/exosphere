package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
)

// FetchTemplatesFeatureContext defines the festure context for features/fetch_templates.feature
// nolint gocyclo
func FetchTemplatesFeatureContext(s *godog.Suite) {
	s.Step(`^I am in the root directory of an empty git application repository called "([^"]*)"$`, func(appName string) error {
		appDir = path.Join(os.TempDir(), appName)
		createEmptyApp(appName, cwd)
		if _, _, _, err := processHelpers.Start("git init", appDir); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to creates .git for %s:%s\n", appName, err))
		}
		return nil
	})

	s.Step(`^my application contains the file "([^"]*)" with the content:$`, func(fileName string, content *gherkin.DocString) error {
		return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(content.Content), 0777)
	})

	s.Step(`^I update "([^"]*)" with the content:$`, func(fileName string, content *gherkin.DocString) error {
		return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(content.Content), 0777)
	})

	s.Step(`^my application contains the directory "([^"]*)"`, func(directory string) error {
		f, err := os.Stat(path.Join(appDir, directory))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to get information for the entry %s", directory))
		}
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", directory)
		} else if !f.IsDir() {
			return fmt.Errorf("%s is a not a directory", directory)
		}
		if osHelpers.IsEmpty(directory) {
			return fmt.Errorf("%s is empty", directory)
		}
		return nil
	})

	s.Step(`^my git repository has "([^"]*)" as a submodule$`, func(submodulePath string) error {
		if _, _, stdoutBuffer, err := processHelpers.Start("git config --file .gitmodules --get-regexp path", appDir); err != nil {
			return errors.Wrap(err, "Failed to read git config file")
		} else {
			return waitForText(stdoutBuffer, submodulePath, 1000)
		}
	})

	s.Step(`^the git URL of "([^"]*)" is "([^"]*)"$`, func(submodulePath, gitURL string) error {
		if _, _, stdoutBuffer, err := processHelpers.Start("git config --file .gitmodules --get-regexp url", appDir); err != nil {
			return errors.Wrap(err, "Failed to read git config file")
		} else {
			return waitForText(stdoutBuffer, gitURL, 1000)
		}
	})

}
