package main

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-run-go/app_runner"
	"github.com/Originate/exosphere/exo-run-go/helpers"
	"github.com/Originate/exosphere/exo-run-go/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo run",
	Short: "Create a new Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		}
		appConfig := helpers.GetAppConfig()
		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		services, silencedServices := helpers.GetExistingServices(appConfig)
		silencedDependencies := helpers.GetSilencedDependencies(appConfig)
		logger := logger.NewLogger(services, append(silencedServices, silencedDependencies...))
		appRunner := appRunner.NewAppRunner(appConfig, logger)
		appRunner.Start()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
