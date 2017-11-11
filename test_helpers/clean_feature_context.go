package testHelpers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/moby/moby/client"
)

func addFile(appDir, serviceFolder, fileName string) error {
	filePath := path.Join(appDir, serviceFolder, fileName)
	return ioutil.WriteFile(filePath, []byte("test"), 0644)
}

// CleanFeatureContext defines the festure context for features/clean.feature
// nolint gocyclo
func CleanFeatureContext(s *godog.Suite) {
	var dockerClient *client.Client
	var appNetwork = "cleancontainers"
	var testNetwork = "cleancontainerstests"
	var thirdPartyContainer = "exosphere-third-party-test-container"
	var appContainerProcess *execplus.CmdPlus
	var testContainerProcess *execplus.CmdPlus
	var thirdPartyContainerProcess *execplus.CmdPlus

	s.BeforeSuite(func() {
		var err error
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			panic(err)
		}
	})

	s.AfterScenario(func(arg1 interface{}, arg2 error) {
		if appContainerProcess != nil {
			err := appContainerProcess.Kill()
			if err != nil {
				panic(err)
			}
			_, err = util.Run(path.Join(appDir, "docker-compose"), "docker-compose", "-p", appNetwork, "-f", "run_development.yml", "down")
			if err != nil {
				panic(err)
			}
		}
		appContainerProcess = nil
		if testContainerProcess != nil {
			err := testContainerProcess.Kill()
			if err != nil {
				panic(err)
			}
			_, err = util.Run(path.Join(appDir, "docker-compose"), "docker-compose", "-p", testNetwork, "-f", "test.yml", "down")
			if err != nil {
				panic(err)
			}
		}
		testContainerProcess = nil
		if thirdPartyContainerProcess != nil {
			err := thirdPartyContainerProcess.Kill()
			if err != nil {
				panic(err)
			}
			timeout := time.Second * 15
			err = dockerClient.ContainerStop(context.Background(), thirdPartyContainer, &timeout)
			if err != nil {
				panic(err)
			}
		}
		thirdPartyContainerProcess = nil
	})

	s.Step(`^my machine has both dangling and non-dangling Docker images and volumes$`, func() error {
		tempAppDir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		serviceRole := "mongo"
		err = CheckoutApp(tempAppDir, "external-dependency")
		if err != nil {
			return fmt.Errorf("Error checking out app: %v", err)
		}
		err = runApp(tempAppDir, "MongoDB connected")
		if err != nil {
			return fmt.Errorf("Error setting up app (first time): %v", err)
		}
		err = addFile(tempAppDir, serviceRole, "test.txt")
		if err != nil {
			return fmt.Errorf("Error adding file: %v", err)
		}
		return killTestContainers(tempAppDir)
	})

	s.Step(`^my machine has running application and test containers$`, func() error {
		var err error
		appContainerProcess, err = runComposeInNetwork("up", appNetwork, path.Join(appDir, "docker-compose"), "run_development.yml")
		if err != nil {
			return err
		}
		err = appContainerProcess.WaitForText("Creating app-test-container", time.Second*5)
		if err != nil {
			return err
		}
		testContainerProcess, err = runComposeInNetwork("up", testNetwork, path.Join(appDir, "docker-compose"), "test.yml")
		if err != nil {
			return err
		}
		return testContainerProcess.WaitForText("Creating service-test-container", time.Second*5)
	})

	s.Step(`^my machine has stopped application and test containers$`, func() error {
		var err error
		appContainerProcess, err = runComposeInNetwork("create", appNetwork, path.Join(appDir, "docker-compose"), "run_development.yml")
		if err != nil {
			return err
		}
		err = appContainerProcess.WaitForText("Creating app-test-container", time.Second*5)
		if err != nil {
			return err
		}
		testContainerProcess, err = runComposeInNetwork("create", testNetwork, path.Join(appDir, "docker-compose"), "test.yml")
		if err != nil {
			return err
		}
		return testContainerProcess.WaitForText("Creating service-test-container", time.Second*20)
	})

	s.Step(`^my machine has running third party containers$`, func() error {
		thirdPartyContainerProcess = execplus.NewCmdPlus("docker", "run", "--name", thirdPartyContainer, "--rm", "alpine", "sleep", "1000")
		err := thirdPartyContainerProcess.Start()
		if err != nil {
			return err
		}
		return waitForContainer(dockerClient, thirdPartyContainer)
	})

	s.Step(`^it removes application and test containers$`, func() error {
		doesHaveAppContainer, err := hasContainer(dockerClient, "app-test-container")
		if err != nil {
			return err
		}
		if doesHaveAppContainer {
			return errors.New("expected app-test-container to be removed but it wasn't")
		}
		doesHaveServiceContainer, err := hasContainer(dockerClient, "service-test-container")
		if err != nil {
			return err
		}
		if doesHaveServiceContainer {
			return errors.New("expected service-test-container to be removed but it wasn't")
		}
		hasAppNetwork, err := hasNetwork(dockerClient, appNetwork)
		if err != nil {
			return err
		}
		if hasAppNetwork {
			return fmt.Errorf("expected network '%s' to have been removed", appNetwork)
		}
		hasTestNetwork, err := hasNetwork(dockerClient, testNetwork)
		if err != nil {
			return err
		}
		if hasTestNetwork {
			return fmt.Errorf("expected network '%s' to have been removed", testNetwork)
		}
		return nil
	})

	s.Step(`^it does not stop any third party containers$`, func() error {
		containers, err := listContainersInNetwork(dockerClient, "bridge")
		if err != nil {
			return err
		}
		hasContainer := false
		for _, container := range containers {
			if container == thirdPartyContainer {
				hasContainer = true
			}
		}
		if !hasContainer {
			return fmt.Errorf("expected third party container '%s' to be running", thirdPartyContainer)
		}
		return nil
	})

	s.Step(`^it has non-dangling images$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "false")
		imageSummaries, err := dockerClient.ImageList(ctx, types.ImageListOptions{
			All:     false,
			Filters: filtersArgs,
		})
		if err != nil {
			return err
		}
		if len(imageSummaries) == 0 {
			return fmt.Errorf("expected non-dangling images but there are none")
		}
		return nil
	})

	s.Step(`^it does not have dangling images$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "true")
		imageSummaries, err := dockerClient.ImageList(ctx, types.ImageListOptions{
			All:     false,
			Filters: filtersArgs,
		})
		if err != nil {
			return err
		}
		if len(imageSummaries) != 0 {
			return fmt.Errorf("expected no dangling images but there are %d", len(imageSummaries))
		}
		return nil
	})

	s.Step(`^it does not have dangling volumes$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "true")
		volumesListOKBody, err := dockerClient.VolumeList(ctx, filtersArgs)
		if err != nil {
			return err
		}
		if len(volumesListOKBody.Volumes) != 0 {
			return fmt.Errorf("expected no dangling volumes but there are %d", len(volumesListOKBody.Volumes))
		}
		return nil
	})
}
