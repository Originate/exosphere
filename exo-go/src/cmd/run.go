package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Originate/exosphere/exo-go/src/application"
	"github.com/Originate/exosphere/exo-go/src/applicationrunner"
	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/osplus"
	"github.com/pkg/errors"
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
		homeDir, err := osplus.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		appConfig, err := config.NewAppConfig(appDir)
		if err != nil {
			panic(err)
		}
		serviceNames := appConfig.GetServiceNames()
		dependencyNames := appConfig.GetDependencyNames()
		silencedServiceNames := appConfig.GetSilencedServiceNames()
		silencedDependencyNames := appConfig.GetSilencedDependencyNames()
		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-run")
		logger := logger.New(roles, append(silencedServiceNames, silencedDependencyNames...), os.Stdout)

		fmt.Printf("Setting up %s %s\n\n", appConfig.Name, appConfig.Version)
		initializer, err := application.NewInitializer(appConfig, logger, appDir, homeDir)
		if err != nil {
			panic(err)
		}
		err = initializer.Initialize()
		if err != nil {
			panic(errors.Wrap(err, "setup failed"))
		}
		fmt.Println("setup complete")

		fmt.Printf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		runner, err := applicationrunner.NewRunner(appConfig, logger, appDir, homeDir)
		if err != nil {
			panic(err)
		}
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			signal.Stop(c)
			if err := runner.Shutdown(config.ShutdownConfig{CloseMessage: " shutting down ..."}); err != nil {
				panic(err)
			}
			wg.Done()
		}()
		if err := runner.Start(); err != nil {
			errorMessage := fmt.Sprint(err)
			if err := runner.Shutdown(config.ShutdownConfig{ErrorMessage: errorMessage}); err != nil {
				panic(err)
			}
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
