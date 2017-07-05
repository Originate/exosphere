package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
)

// CreateFeatureContext defines the festure context for features/create.feature
// nolint gocyclo
func CreateFeatureContext(s *godog.Suite) {
	s.Step(`^my workspace contains the file "([^"]*)" with content:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})

	s.Step(`^my workspace contains the empty directory "([^"]*)"`, func(directory string) error {
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
