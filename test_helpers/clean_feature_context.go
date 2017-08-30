package testHelpers

import (
	"context"
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

func addFile(cwd, appName, serviceFolder, fileName string) error {
	filePath := path.Join(cwd, "tmp", appName, serviceFolder, fileName)
	return ioutil.WriteFile(filePath, []byte("test"), 0644)
}

// CleanFeatureContext defines the festure context for features/clean.feature
// nolint gocyclo
func CleanFeatureContext(s *godog.Suite) {
	var exosphereAppNetwork = "exosphereapptesting"
	var exosphereTestNetwork = "exospheretesttesting"
	var thirdPartyContainer = "exosphere-third-party-test-container"
	var dockerClient *client.Client
	var appContainerProcess *execplus.CmdPlus
	var serviceTestContainerProcess *execplus.CmdPlus
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
			_, err = util.Run(path.Join(appDir, "tmp"), "docker-compose", "-p", exosphereAppNetwork, "down")
			if err != nil {
				panic(err)
			}
		}
		appContainerProcess = nil
		if serviceTestContainerProcess != nil {
			err := serviceTestContainerProcess.Kill()
			if err != nil {
				panic(err)
			}
			_, err = util.Run(path.Join(appDir, "service", "tests", "tmp"), "docker-compose", "-p", exosphereTestNetwork, "down")
			if err != nil {
				panic(err)
			}
		}
		serviceTestContainerProcess = nil
		if thirdPartyContainerProcess != nil {
			err := thirdPartyContainerProcess.Kill()
			if err != nil {
				panic(err)
			}
		}
		thirdPartyContainerProcess = nil
	})

	s.Step(`^my machine has both dangling and non-dangling Docker images and volumes$`, func() error {
		appName := "external-dependency"
		serviceName := "mongo"
		err := CheckoutApp(cwd, appName)
		if err != nil {
			return fmt.Errorf("Error checking out app: %v", err)
		}
		err = runApp(cwd, appName)
		if err != nil {
			return fmt.Errorf("Error setting up app (first time): %v", err)
		}
		err = addFile(cwd, appName, serviceName, "test.txt")
		if err != nil {
			return fmt.Errorf("Error adding file: %v", err)
		}
		dockerComposeDir := path.Join(appDir, "tmp")
		return killTestContainers(dockerComposeDir)
	})

	s.Step(`^my machine has running application and test containers$`, func() error {
		appContainerProcess = execplus.NewCmdPlus("docker-compose", "-p", exosphereAppNetwork, "up")
		appContainerProcess.SetDir(path.Join(appDir, "tmp"))
		err := appContainerProcess.Start()
		if err != nil {
			return err
		}
		serviceTestContainerProcess = execplus.NewCmdPlus("docker-compose", "-p", exosphereTestNetwork, "up")
		serviceTestContainerProcess.SetDir(path.Join(appDir, "service", "tests", "tmp"))
		err = serviceTestContainerProcess.Start()
		time.Sleep(2500 * time.Millisecond)
		return err
	})

	s.Step(`^my machine has stopped application and test containers$`, func() error {
		appContainerProcess = execplus.NewCmdPlus("docker-compose", "-p", exosphereAppNetwork, "create")
		appContainerProcess.SetDir(path.Join(appDir, "tmp"))
		err := appContainerProcess.Start()
		if err != nil {
			return err
		}
		serviceTestContainerProcess = execplus.NewCmdPlus("docker-compose", "-p", exosphereTestNetwork, "create")
		serviceTestContainerProcess.SetDir(path.Join(appDir, "service", "tests", "tmp"))
		return serviceTestContainerProcess.Start()
	})

	s.Step(`^my machine has running third party containers$`, func() error {
		serviceTestContainerProcess = execplus.NewCmdPlus("docker", "run", "--name", thirdPartyContainer, "--rm", "alpine", "sleep", "20")
		time.Sleep(2500 * time.Millisecond)
		return serviceTestContainerProcess.Start()
	})

	s.Step(`^it removes application and test containers$`, func() error {
		hasNetworkBool, err := hasNetwork(dockerClient, exosphereAppNetwork)
		if err != nil {
			return err
		}
		if hasNetworkBool {
			return fmt.Errorf("Expected network '%s' to have been removed.", exosphereAppNetwork)
		}
		hasNetworkBool, err = hasNetwork(dockerClient, exosphereTestNetwork)
		if err != nil {
			return err
		}
		if hasNetworkBool {
			return fmt.Errorf("Expected network '%s' to have been removed.", exosphereTestNetwork)
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
			return fmt.Errorf("Expected third party container '%s' to be running.", thirdPartyContainer)
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
			return fmt.Errorf("Expected non-dangling images but there are none")
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
			return fmt.Errorf("Expected no dangling images but there are %d", len(imageSummaries))
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
			return fmt.Errorf("Expected no dangling volumes but there are %d", len(volumesListOKBody.Volumes))
		}
		return nil
	})
}
