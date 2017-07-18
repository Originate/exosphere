package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes dangling Docker images and volumes",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to clean up your Docker workspace!\n\n")
		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		err = dockerHelpers.RemoveDanglingImages(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling images")
		err = dockerHelpers.RemoveDanglingVolumes(c)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling volumes")
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
