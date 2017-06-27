package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kr/pretty"
	"github.com/mattn/go-shellwords"
)

// the temp dir in which the test repos live
var testRoot string

// the output of the last command run
var output string

// Cucumber step definitions
func FeatureContext(s *godog.Suite) {

	// the error of the last run operation
	var err error

	s.Step(`^it prints:$`, func(text *gherkin.DocString) error {
		matched, err := regexp.Match(text.Content, []byte(output))
		check(err)
		if !matched {
			return errors.New("output does not match")
		}
		return nil
	})

	s.Step(`^running "([^"]*)"$`, func(command string) (err error) {
		command = makeCrossPlatformCommand(command)
		words, err := shellwords.Parse(command)
		check(err)
		err = run(words)
		fmt.Println(output)
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
