package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_runner"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an Exosphere application",
	Long:  "Runs an Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		appConfig := appConfigHelpers.GetAppConfig()
		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)

		serviceNames := appConfigHelpers.GetServiceNames(appConfig.Services)
		dependencyNames := appConfigHelpers.GetDependencyNames(appConfig)
		silencedServiceNames := appConfigHelpers.GetSilencedServiceNames(appConfig.Services)
		silencedDependencyNames := appConfigHelpers.GetSilencedDependencyNames(appConfig)

		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-run")
		logger := logger.NewLogger(roles, append(silencedServiceNames, silencedDependencyNames...))
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get the current path: %s", err)
		}
		appRunner := appRunner.NewAppRunner(appConfig, logger, cwd)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			if output, err := appRunner.Shutdown(" shutting down ...", ""); err != nil {
				log.Fatalf("Failed to shutdown the app\nOutput: %s\nError: %s\n", output, err)
			}
			wg.Done()
		}()
		if output, err := appRunner.Start(); err != nil {
			errorMessage := fmt.Sprintf("Failed to run images\nOutput: %s\nErrorL %s\n", output, err)
			if output, err := appRunner.Shutdown("", errorMessage); err != nil {
				log.Fatalf("Failed to shutdown the app\nOutput: %s\nError: %s\n", output, err)
			}
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
