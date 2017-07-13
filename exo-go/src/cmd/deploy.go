package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
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
		terraformFileHelpers.GenerateTerraform(appConfig)
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
