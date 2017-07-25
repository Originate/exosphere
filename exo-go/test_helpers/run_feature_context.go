package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// RunFeatureContext defines the festure context for features/run
// nolint gocyclo
func RunFeatureContext(s *godog.Suite) {
	var dockerClient *client.Client
	var cwd, appName string

	s.BeforeSuite(func() {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^I am in the root directory of the "([^"]*)" example application$`, func(name string) error {
		appDir = path.Join(cwd, "tmp", name)
		appName = name
		return CheckoutApp(cwd, appName)
	})

	s.Step(`^my application has been set up correctly$`, func() error {
		if err := setupApp(cwd, appName); err != nil {
			return err
		}
		return nil
	})

	s.Step(`^my machine is running the services:$`, func(table *gherkin.DataTable) error {
		runningContainers, err := dockerHelpers.ListRunningContainers(dockerClient)
		if err != nil {
			return errors.Wrap(err, "Failed to list running containers")
		}
		for _, row := range table.Rows[1:] {
			serviceName := row.Cells[0].Value
			if !util.DoesStringArrayContain(runningContainers, serviceName) {
				err = fmt.Errorf("Expected the machine to be running the service '%s'", serviceName)
				break
			}
		}
		return err
	})

	s.Step(`^adding a file to "([^"]*)" service folder$`, func(serviceDir string) error {
		return ioutil.WriteFile(path.Join(appDir, serviceDir, "test.txt"), []byte(""), 0777)
	})

	s.Step(`^the "([^"]*)" service restarts$`, func(serviceName string) error {
		if err := process.WaitForTextWithTimeout(fmt.Sprintf("Restarting service '%s'", serviceName), 5000); err != nil {
			return err
		}
		return process.WaitForTextWithTimeout(fmt.Sprintf("'%s' restarted successfully", serviceName), 5000)
	})

}
