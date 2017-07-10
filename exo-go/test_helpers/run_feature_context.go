package testHelpers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
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

	s.Step(`^my machine is running ExoCom$`, func() error {
		return waitForText(stdoutBuffer, "ExoCom WebSocket listener online at port", 5000)
	})

	s.Step(`^ExoCom uses this routing:$`, func(table *gherkin.DataTable) error {
		expectedRoutes := []types.ServiceRoute{}
		for _, row := range table.Rows[1:] {
			role, receives, sends, namespace := row.Cells[0].Value, row.Cells[1].Value, row.Cells[2].Value, row.Cells[3].Value
			expectedRoutes = append(expectedRoutes, types.ServiceRoute{role, strings.Split(receives, ", "), strings.Split(sends, ", "), namespace})
		}
		dockerCompose := getDockerCompose()
		actualRoutes := dockerCompose.Services["exocom0.22.1"].Environment.ServiceRoutes
		fmt.Println(expectedRoutes)
		fmt.Println(actualRoutes)
		fmt.Println(reflect.DeepEqual(expectedRoutes, actualRoutes))
		// if !reflect.DeepEqual(expectedRoutes, actualRoutes) {
		// 	return fmt.Errorf("Expected exocom to use this routing")
		// }
		return nil
	})

	s.Step(`^the web service broadcasts a "([^"]*)" message$`, func(message string) error {
		_, err := http.Get("http://localhost:4000")
		return err
	})

	s.Step(`^the "([^"]*)" service receives a "([^"]*)" message$`, func(service, message string) error {
		return waitForText(stdoutBuffer, fmt.Sprintf("'%s' service received message '%s'", service, message), 5000)
	})

}
