package main_test

import (
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exosphere/test/helpers"
)

func TestMain(m *testing.M) {
	format := "pretty"
	var paths []string
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "features/") {
			paths = append(paths, arg)
		}
	}
	if len(paths) == 0 {
		format = "progress"
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		helpers.AddFeatureContext(s)
		helpers.CleanFeatureContext(s)
		helpers.TemplateFeatureContext(s)
		helpers.SharedFeatureContext(s)
		helpers.RunFeatureContext(s)
		helpers.TestFeatureContext(s)
		helpers.TutorialFeatureContext(s)
	}, godog.Options{
		Format:        format,
		NoColors:      false,
		StopOnFailure: true,
		Paths:         paths,
		Strict:        true,
	})

	os.Exit(status)
}
