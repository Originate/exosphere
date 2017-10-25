package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/util"
	"github.com/docker/docker/api/types/filters"
	"github.com/moby/moby/client"
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
		writer := os.Stdout
		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		util.PrintHeader("Removing dangling images")
		_, err = c.ImagesPrune(context.Background(), filters.NewArgs())
		if err != nil {
			panic(err)
		}
		util.PrintHeader("Removing dangling volumes")
		_, err = c.VolumesPrune(context.Background(), filters.NewArgs())
		if err != nil {
			panic(err)
		}
		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		util.PrintHeader("Removing application containers")
		composeProjectName := composebuilder.GetDockerComposeProjectName(appDir)
		err = application.CleanApplicationContainers(appDir, composeProjectName, writer)
		if err != nil {
			panic(err)
		}
		util.PrintHeader("Removing test containers")
		testDockerComposeProjectName := getTestDockerComposeProjectName(appDir)
		err = application.CleanServiceTestContainers(appDir, testDockerComposeProjectName, writer)
		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
