package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
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

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %s", err)
		}

		appConfig, err := appConfigHelpers.GetAppConfig(cwd)
		if err != nil {
			log.Fatalf("Cannot read application configuration: %s", err)
		}

		serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(cwd, appConfig)
		if err != nil {
			log.Fatalf("Failed to read service configurations: %s", err)
		}

		err = terraformFileHelpers.GenerateTerraformFile(appConfig, serviceConfigs, cwd)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
