package main

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/mattn/go-shellwords"
)

// the output of the last command run
var output string

// Cucumber step definitions
func FeatureContext(s *godog.Suite) {

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
		words, err := shellwords.Parse(command)
		if err != nil {
			return err
		}
		output, err = run(words)
		if err != nil {
			return err
		}
		return
	})

}

func TestMain(m *testing.M) {
	var paths []string
	if len(os.Args) == 2 {
		paths = append(paths, strings.Split(os.Args[1], "=")[1])
	} else {
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:        "pretty",
		NoColors:      false,
		StopOnFailure: true,
		Paths:         paths,
	})

	os.Exit(status)
}
