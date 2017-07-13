package cmd

import (
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys Exosphere application to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about to deploy an application!")

		appConfig := appConfigHelpers.GetAppConfig()
		serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appConfig)
		if err != nil {
			log.Fatalf("Error compiling service configurations: %s", err)
		}
		terraformFileHelpers.GenerateTerraform(appConfig, serviceConfigs)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
