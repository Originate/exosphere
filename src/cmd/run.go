package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var productionFlag bool

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
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		appConfig, err := types.NewAppConfig(appDir)
		if err != nil {
			panic(err)
		}
		dockerComposeProjectName := composebuilder.GetDockerComposeProjectName(appDir)
		serviceNames := appConfig.GetServiceNames()
		dependencyNames := appConfig.GetDependencyNames()
		silencedServiceNames := appConfig.GetSilencedServiceNames()
		silencedDependencyNames := appConfig.GetSilencedDependencyNames()
		logRole := "exo-run"
		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, logRole)
		logger := application.NewLogger(roles, append(silencedServiceNames, silencedDependencyNames...), os.Stdout)
		logChannel := logger.GetLogChannel("exo-run")

		logChannel <- fmt.Sprintf("Setting up %s %s\n\n", appConfig.Name, appConfig.Version)
		initializer, err := application.NewInitializer(
			appConfig,
			logChannel,
			logRole,
			appDir,
			homeDir,
			dockerComposeProjectName,
			productionFlag,
		)
		if err != nil {
			panic(err)
		}
		err = initializer.Initialize()
		if err != nil {
			panic(errors.Wrap(err, "setup failed"))
		}
		logChannel <- "setup complete"

		logChannel <- fmt.Sprintf("Running %s %s\n\n", appConfig.Name, appConfig.Version)
		runner, err := application.NewRunner(appConfig, logger, logRole, appDir, homeDir, dockerComposeProjectName)
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
	runCmd.PersistentFlags().BoolVarP(&productionFlag, "production", "", false, "Run in production mode")
}
