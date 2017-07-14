package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about to deploy an application!")

		appConfig, err := appConfigHelpers.GetAppConfig()
		if err != nil {
			log.Fatalf("Failed to read application configuration: %s", err)
		}

		serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appConfig)
		if err != nil {
			log.Fatalf("Failed to read service configurations: %s", err)
		}

		err = terraformFileHelpers.GenerateTerraform(appConfig, serviceConfigs)
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
