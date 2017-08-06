package docker

import (
	"context"
	"io/ioutil"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

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
