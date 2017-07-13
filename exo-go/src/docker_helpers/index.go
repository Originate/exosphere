package dockerHelpers

import (
	"context"

	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CatFileInDockerImage reads the file fileName inside the given image
func CatFileInDockerImage(c *client.Client, image, fileName string) ([]byte, error) {
	if err := PullImage(c, image); err != nil {
		return []byte(""), err
	}
	output, err := processHelpers.Run("", "docker", "run", image, "cat", fileName)
	return []byte(output), err
}

// PullImage pulls the given image from DockerHub, returns an error if any
func PullImage(c *client.Client, image string) error {
	ctx := context.Background()
	_, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	return err
}
