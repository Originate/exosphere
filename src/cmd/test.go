package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/Originate/exosphere/src/application"
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
			Mount:       true,
			Environment: composebuilder.BuildModeEnvironmentTest,
		}
		if noMountFlag {
			buildMode.Mount = false
		}

		// Listen to SIGINT
		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)

		var testResult types.TestResult
		if context.HasServiceContext {
			testResult, err = application.TestService(context.ServiceContext, writer, buildMode, shutdownChannel)
		} else {
			testResult, err = application.TestApp(context.AppContext, writer, buildMode, shutdownChannel)
		}
		if err != nil {
			panic(err)
		}
		signal.Stop(shutdownChannel)
		if !testResult.Passed && !testResult.Interrupted {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run tests without mounting")
}
