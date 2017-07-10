package dockerHelpers

import (
	"context"
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// RemoveDanglingImages removes all dangling images on the machine
func RemoveDanglingImages(c *client.Client) error {
	ctx := context.Background()
	filtersArgs := filters.NewArgs()
	filtersArgs.Add("dangling", "true")
	imageSummaries, err := c.ImageList(ctx, types.ImageListOptions{
		All:     false,
		Filters: filtersArgs,
	})
	if err != nil {
		return err
	}
	for _, imageSummary := range imageSummaries {
		_, err = c.ImageRemove(ctx, imageSummary.ID, types.ImageRemoveOptions{})
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

// PullImage pulls the given image from DockerHub
func PullImage(c *client.Client, image string, done func(error)) {
	ctx := context.Background()
	_, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	done(err)
}

// CatFile reads the file fileName inside the docker image image
func CatFile(c *client.Client, image, fileName string, done func(err error, content []byte)) {
	PullImage(c, image, func(err error) {
		if err != nil {
			done(err, []byte(""))
		} else {
			output, err := processHelpers.Run(fmt.Sprintf("docker run %s cat %s", image, fileName))
			done(err, []byte(output))
		}
	})
}

// ListRunningContainers passes a slice of the names of running containers
// and error (if any) to the callback function
func ListRunningContainers(c *client.Client, done func([]string, error)) {
	containerNames := []string{}
	ctx := context.Background()
	containers, err := c.ContainerList(ctx, types.ContainerListOptions{})
	for _, container := range containers {
		containerNames = append(containerNames, container.Image)
	}
	done(containerNames, err)
}
