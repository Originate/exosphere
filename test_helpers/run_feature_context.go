package testHelpers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/util"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
)

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

	s.Step("^the docker images have the following folders:", func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			imageName, folder := row.Cells[0].Value, row.Cells[1].Value
			content, err := tools.RunInDockerImage(imageName, "ls")
			if err != nil {
				return err
			}
			folders := strings.Split(content, "\n")
			if !util.DoesStringArrayContain(folders, folder) {
				return fmt.Errorf("Expected the docker image '%s' to have the folder '%s'", imageName, folder)
			}
		}
		return nil
	})

	s.Step(`^my machine has acquired the Docker images:$`, func(table *gherkin.DataTable) error {
		images, err := tools.ListImages(dockerClient)
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
		runningContainers, err := tools.ListRunningContainers(dockerClient)
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

	s.Step(`^modifying (\S*) to "([^"]*)"$`, func(filePath string, fileContent string) error {
		return ioutil.WriteFile(path.Join(appDir, filePath), []byte(fileContent), 0777)
	})

	s.Step(`^my machine contains the network "([^"]*)"`, func(networkName string) error {
		hasNetwork, err := hasNetwork(dockerClient, networkName)
		if err != nil {
			return err
		}
		if !hasNetwork {
			err = fmt.Errorf("Expected the machine to have network %s", networkName)
		}
		return err
	})

	s.Step(`^the network "([^"]*)" contains the running services:$`, func(networkName string, table *gherkin.DataTable) error {
		runningContainers, err := listContainersInNetwork(dockerClient, networkName)
		if err != nil {
			return err
		}
		for _, row := range table.Rows[1:] {
			serviceName := row.Cells[0].Value
			if !util.DoesStringArrayContain(runningContainers, serviceName) {
				err = fmt.Errorf("Expected network to be running the service '%s'", serviceName)
				break
			}
		}
		return err
	})
}
