package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
)

// CatFileInDockerImage reads the file fileName inside the given image
func CatFileInDockerImage(c *client.Client, imageName, fileName string) ([]byte, error) {
	if err := PullImage(c, imageName); err != nil {
		return []byte(""), err
	}
	command := fmt.Sprintf("cat %s", fileName)
	output, err := RunInDockerImage(imageName, command)
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

// ListRunningContainers returns the names of running containers
// and an error (if any)
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

// ListImages returns the names of all images and an error (if any)
func ListImages(c *client.Client) ([]string, error) {
	imageNames := []string{}
	ctx := context.Background()
	imageSummaries, err := c.ImageList(ctx, dockerTypes.ImageListOptions{
		All: true,
	})
	if err != nil {
		return imageNames, err
	}
	for _, imageSummary := range imageSummaries {
		if len(imageSummary.RepoTags) > 0 {
			repoTag := imageSummary.RepoTags[0]
			imageNames = append(imageNames, strings.Split(repoTag, ":")[0])
		}
	}
	return imageNames, nil
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

// RunInDockerImage runs the given command in a new writeable container layer
// over the given image, removes the container when the command exits, and returns
// the output string and an error if any
func RunInDockerImage(imageName, command string) (string, error) {
	return util.Run("", fmt.Sprintf("docker run --rm %s %s", imageName, command))
}
