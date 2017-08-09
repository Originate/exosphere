package testHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo/src/docker"
	"github.com/Originate/exosphere/exo/src/util"
	"github.com/moby/moby/client"
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

	s.Step("^the docker images have the following folders:", func(table *gherkin.DataTable) error {
		var err error
		for _, row := range table.Rows[1:] {
			imageName, folder := row.Cells[0].Value, row.Cells[1].Value
			content, err := docker.RunInDockerImage(imageName, "ls")
			if err != nil {
				return err
			}
			folders := strings.Split(content, "\n")
			if !util.DoesStringArrayContain(folders, folder) {
				err = fmt.Errorf("Expected the docker image '%s' to have the folder", imageName, folder)
				break
			}
		}
		return err
	})

	s.Step(`^my machine has acquired the Docker images:$`, func(table *gherkin.DataTable) error {
		images, err := docker.ListImages(dockerClient)
		if err != nil {
			return errors.Wrap(err, "Failed to list docker images")
		}
		for _, row := range table.Rows[1:] {
			imageName := row.Cells[0].Value
			if !util.DoesStringArrayContain(images, imageName) {
				err = fmt.Errorf("Expected the machine to have acquired the docker image '%s'", imageName)
				break
			}
		}
		return err
	})

	s.Step(`^my machine is running the services:$`, func(table *gherkin.DataTable) error {
		runningContainers, err := docker.ListRunningContainers(dockerClient)
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
		if err := childCmdPlus.WaitForText(fmt.Sprintf("Restarting service '%s'", serviceName), time.Second*5); err != nil {
			return err
		}
		return childCmdPlus.WaitForText(fmt.Sprintf("'%s' restarted successfully", serviceName), time.Second*5)
	})

}
