package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo clean",
	Short: "Removes dangling Docker images and volumes.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to clean up your Docker workspace!\n\n")
		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		err = removeImages(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling images")
		err = removeVolumes(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling volumes")
	},
}

func removeImages(c *client.Client) error {
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

func removeVolumes(c *client.Client) error {
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

// nolint unused
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
