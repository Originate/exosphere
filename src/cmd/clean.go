package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
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
		_, err = c.ImagesPrune(context.Background(), filters.NewArgs())
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling images")
		_, err = c.VolumesPrune(context.Background(), filters.NewArgs())
		if err != nil {
			panic(err)
		}
		fmt.Println("removed all dangling volumes")
		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		composeProjectName := composebuilder.GetDockerComposeProjectName(appDir)
		err = application.CleanApplicationContainers(appDir, composeProjectName, writer)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed application containers")
		testDockerComposeProjectName := getTestDockerComposeProjectName(appDir)
		err = application.CleanServiceTestContainers(appDir, testDockerComposeProjectName, writer)
		if err != nil {
			panic(err)
		}
		fmt.Println("removed test containers")
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
