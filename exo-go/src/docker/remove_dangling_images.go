package docker

import (
	"context"

	dockerTypes "github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/filters"

	"github.com/moby/moby/client"
)

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
