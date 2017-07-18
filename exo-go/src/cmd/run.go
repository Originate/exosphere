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
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an Exosphere application",
	Long:  "Runs an Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get the current path: %s", err)
		}
		appConfig, err := appConfigHelpers.GetAppConfig(cwd)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		serviceNames := appConfigHelpers.GetServiceNames(appConfig.Services)
		dependencyNames := appConfigHelpers.GetDependencyNames(appConfig)
		silencedServiceNames := appConfigHelpers.GetSilencedServiceNames(appConfig.Services)
		silencedDependencyNames := appConfigHelpers.GetSilencedDependencyNames(appConfig)

		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-run")
		logger := logger.NewLogger(roles, append(silencedServiceNames, silencedDependencyNames...))

		appRunner := appRunner.NewAppRunner(appConfig, logger, cwd)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			signal.Stop(c)
			if err := appRunner.Shutdown(types.ShutdownConfig{CloseMessage: " shutting down ..."}); err != nil {
				panic(err)
			}
			wg.Done()
		}()
		if err := appRunner.Start(); err != nil {
			errorMessage := fmt.Sprint(err)
			if err := appRunner.Shutdown(types.ShutdownConfig{ErrorMessage: errorMessage}); err != nil {
				panic(err)
			}
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
