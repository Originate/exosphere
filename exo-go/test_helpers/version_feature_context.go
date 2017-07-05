package testHelpers

import (
	"errors"
	"regexp"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// VersionFeatureContext defines the festure context for features/version.feature
func VersionFeatureContext(s *godog.Suite) {
	var output string

	s.Step(`^it prints:$`, func(expected *gherkin.DocString) error {
		matched, err := regexp.Match(expected.Content, []byte(output))
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("output does not match")
		}
		return nil
	})

	s.Step(`^running "([^"]*)"$`, func(command string) (err error) {
		output, err = run(command)
		if err != nil {
			return err
		}
		return
	})
}
