package dockerHelpers

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// CatFileInDockerImage reads the file fileName inside the given image
func CatFileInDockerImage(c *client.Client, image, fileName string) ([]byte, error) {
	if err := PullImage(c, image); err != nil {
		return []byte(""), err
	}
	output, err := processHelpers.Run("", "docker", "run", image, "cat", fileName)
	return []byte(output), err
}

// GetDockerCompose reads docker-compose.yml at the given path and
// returns the dockerCompose object
func GetDockerCompose(dockerComposePath string) (result types.DockerCompose, err error) {
	yamlFile, err := ioutil.ReadFile(dockerComposePath)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "Failed to unmarshal docker-compose.yml")
	}
	return result, nil
}

// GetRenderedVolumes returns the rendered paths to the given volumes
func GetRenderedVolumes(volumes []string, appName string, role string, homeDir string) ([]string, error) {
	dataPath := path.Join(homeDir, ".exosphere", appName, role, "data")
	renderedVolumes := []string{}
	if err := os.MkdirAll(dataPath, 0777); err != nil { //nolint gas
		return renderedVolumes, errors.Wrap(err, "Failed to create the necessary directories for the volumes")
	}
	for _, volume := range volumes {
		renderedVolumes = append(renderedVolumes, strings.Replace(volume, "{{EXO_DATA_PATH}}", dataPath, -1))
	}
	return renderedVolumes, nil
}

// ListRunningContainers passes a slice of the names of running containers
// and error (if any) to the callback function
func ListRunningContainers(c *client.Client) ([]string, error) {
	containerNames := []string{}
	ctx := context.Background()
	containers, err := c.ContainerList(ctx, dockerTypes.ContainerListOptions{})
	if err != nil {
		return containerNames, err
	}
	for _, container := range containers {
		containerNames = append(containerNames, strings.Replace(container.Names[0], "/", "", -1))
	}
	return containerNames, nil
}

// PullImage pulls the given image from DockerHub, returns an error if any
func PullImage(c *client.Client, image string) error {
	ctx := context.Background()
	stream, err := c.ImagePull(ctx, image, dockerTypes.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(stream)
	return err
}

// RemoveDanglingImages removes all dangling images on the machine
func RemoveDanglingImages(c *client.Client) error {
	ctx := context.Background()
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("dangling", "true")
	imageSummaries, err := c.ImageList(ctx, dockerTypes.ImageListOptions{
		All:     false,
		Filters: filtersArgs,
	})
	if err != nil {
		return err
	}
	for _, imageSummary := range imageSummaries {
		_, err = c.ImageRemove(ctx, imageSummary.ID, dockerTypes.ImageRemoveOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveDanglingVolumes removes all dangling volumes on the machine
func RemoveDanglingVolumes(c *client.Client) error {
	ctx := context.Background()
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("dangling", "true")
	volumesListOKBody, err := c.VolumeList(ctx, filtersArgs)
	if err != nil {
		return err
	}
	for _, volume := range volumesListOKBody.Volumes {
		err = c.VolumeRemove(ctx, volume.Name, false)
		if err != nil {
			return err
		}
	}
	return nil
}
