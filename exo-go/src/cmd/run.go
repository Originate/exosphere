package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_runner"
	"github.com/Originate/exosphere/exo-go/src/app_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
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
			panic(err)
		}
		serviceNames := appConfigHelpers.GetServiceNames(appConfig.Services)
		dependencyNames := appConfigHelpers.GetDependencyNames(appConfig)
		silencedServiceNames := appConfigHelpers.GetSilencedServiceNames(appConfig.Services)
		silencedDependencyNames := appConfigHelpers.GetSilencedDependencyNames(appConfig)
		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-run")
		logger := logger.NewLogger(roles, append(silencedServiceNames, silencedDependencyNames...), os.Stdout)

		fmt.Printf("Setting up %s %s\n\n", appConfig.Name, appConfig.Version)
		setup, err := appSetup.NewAppSetup(appConfig, logger, appDir, homeDir)
		if err != nil {
			panic(err)
		}
		err = setup.Setup()
		if err != nil {
			panic(errors.Wrap(err, "setup failed"))
		}
		fmt.Println("setup complete")

		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		runner := appRunner.NewAppRunner(appConfig, logger, appDir, homeDir)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			signal.Stop(c)
			if err := runner.Shutdown(types.ShutdownConfig{CloseMessage: " shutting down ..."}); err != nil {
				panic(err)
			}
			wg.Done()
		}()
		if err := runner.Start(); err != nil {
			errorMessage := fmt.Sprint(err)
			if err := runner.Shutdown(types.ShutdownConfig{ErrorMessage: errorMessage}); err != nil {
				panic(err)
			}
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
