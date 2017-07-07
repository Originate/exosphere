package cmd

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_runner"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an Exosphere application",
	Long: `Runs an Exosphere application
This command must be run in the root directory of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		appConfig := appConfigHelpers.GetAppConfig()
		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		services := serviceHelpers.GetExistingServices(appConfig.Services)
		silencedServices := serviceHelpers.GetSilencedServices(appConfig.Services)
		silencedDependencies := appConfigHelpers.GetSilencedDependencies(appConfig)
		logger := logger.NewLogger(services, append(silencedServices, silencedDependencies...))
		appRunner := appRunner.NewAppRunner(appConfig, logger)
		appRunner.Start()
		// TODO: capture SIGINT and shutdown properly
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
