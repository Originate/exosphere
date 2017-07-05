package main_test

import (
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exosphere/exo-go/test_helpers"
)

func TestMain(m *testing.M) {
	var format string
	var paths []string
	if len(os.Args) == 3 && os.Args[1] == "--" {
		format = "pretty"
		paths = append(paths, os.Args[2])
	} else {
		format = "pretty"
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		testHelpers.AddFeatureContext(s)
		testHelpers.CleanFeatureContext(s)
		testHelpers.CreateFeatureContext(s)
		testHelpers.SharedFeatureContext(s)
	}, godog.Options{
		Format:        format,
		NoColors:      false,
		StopOnFailure: true,
		Paths:         paths,
	})

	os.Exit(status)
}
