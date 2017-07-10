package testHelpers

import (
	"fmt"
	"os"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	gherkin "github.com/cucumber/gherkin-go"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func checkoutAndRunApp(appName, onlineText string) error {
	cwd, err := os.Getwd()
	if err != nil {
		errors.Wrap(err, "Failed to get the current path")
	}
	if err := checkoutApp(cwd, appName); err != nil {
		return err
	}
	appDir = path.Join(cwd, "tmp", appName)
	if err := setupApp(cwd, appName); err != nil {
		return err
	}
	return runApp(cwd, appName, onlineText)
}

// RunFeatureContext defines the festure context for features/run
// nolint gocyclo
func RunFeatureContext(s *godog.Suite) {
	var dockerClient *client.Client

	s.BeforeSuite(func() {
		var err error
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^a running "([^"]*)" application$`, func(appName string) error {
		return checkoutAndRunApp(appName, "all services online")
	})

	s.Step(`^my machine is running the services:$`, func(table *gherkin.DataTable) error {
		var err error
		dockerHelpers.ListRunningContainers(dockerClient, func(runningContainers []string, err error) {
			if err != nil {
				err = errors.Wrap(err, "Failed to list running containers")
				return
			}
			for _, row := range table.Rows[1:] {
				serviceName := row.Cells[0].Value
				if !util.DoesStringArrayContain(runningContainers, serviceName) {
					err = fmt.Errorf("Expected the machine to be running the service '%s'", serviceName)
					break
				}
			}
		})
		return err
	})
}
