package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_deployer"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
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
		homeDir, err := osHelpers.GetUserHomeDir()
		if err != nil {
			panic(err)
		}
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		if err != nil {
			log.Fatalf("Cannot read application configuration: %s", err)
		}

		serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
		if err != nil {
			log.Fatalf("Failed to read service configurations: %s", err)
		}

		logger := logger.NewLogger([]string{"exo-deploy"}, []string{}, os.Stdout)
		appDeployer := appDeployer.NewAppDeployer(appConfig, serviceConfigs, appDir, homeDir, logger)
		if err := appDeployer.Start(); err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}
