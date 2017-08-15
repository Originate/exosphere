package testHelpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/DATA-DOG/godog"
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
	var dockerClient *client.Client

	s.BeforeSuite(func() {
		var err error
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^my machine has both dangling and non-dangling Docker images and volumes$`, func() error {
		appName := "external-dependency"
		serviceName := "mongo-service"
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
