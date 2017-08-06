package docker

import (
	"context"

	"github.com/moby/moby/api/types/filters"
	"github.com/moby/moby/client"
)

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
