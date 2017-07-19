package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Long:  "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about to deploy an application!")
		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		if err != nil {
			log.Fatalf("Cannot read application configuration: %s", err)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %s", err)
		}
		err = terraformFileHelpers.GenerateTerraform(appConfig, cwd)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
