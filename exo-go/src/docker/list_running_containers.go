package docker

import (
	"context"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"

	"github.com/moby/moby/client"
)

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
