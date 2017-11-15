package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/Originate/exosphere/src/application/tester"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs tests for the application",
	Long:  "Runs tests for the application",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}

		context, err := GetContext()
		if err != nil {
			log.Fatal(err)
		}
		writer := os.Stdout
		buildMode := composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Environment: composebuilder.BuildModeEnvironmentTest,
		}

		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)

		var testResult types.TestResult
		if context.HasServiceContext {
			testResult, err = tester.TestService(context, writer, buildMode, shutdownChannel)
		} else {
			testResult, err = tester.TestApp(context.AppContext, writer, buildMode, shutdownChannel)
		}
		if err != nil {
			panic(err)
		}
		signal.Stop(shutdownChannel)
		if !testResult.Passed {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
