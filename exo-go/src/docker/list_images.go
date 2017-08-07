package docker

import (
	"context"
	"strings"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

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
