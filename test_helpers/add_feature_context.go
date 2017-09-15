package testHelpers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pkg/errors"
)

// AddFeatureContext defines the festure context for features/add/**/*.feature
func AddFeatureContext(s *godog.Suite) {
	s.Step(`^my application now contains the file "([^"]*)" containing the text:$`, func(fileName string, expectedContent *gherkin.DocString) error {
		bytes, err := ioutil.ReadFile(path.Join(appDir, fileName))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read %s", fileName))
		}
		return validateTextContains(strings.TrimSpace(string(bytes)), strings.TrimSpace(expectedContent.Content))
	})
}
