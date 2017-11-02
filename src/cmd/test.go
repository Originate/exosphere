package cmd

import (
	"log"
	"os"
	"os/signal"
	"path"

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
		appContext := context.AppContext
		writer := os.Stdout
		buildMode := composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Mount:       true,
			Environment: composebuilder.BuildModeEnvironmentTest,
		}
		if noMountFlag {
			buildMode.Mount = false
		}

		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)

		var testResult types.TestResult
		if context.HasServiceContext {
			var testRunner *tester.TestRunner
			testRunner, err = tester.NewTestRunner(appContext, writer, buildMode)
			if err != nil {
				panic(err)
			}
			testRole := path.Base(context.ServiceContext.Location)
			testResult, err = tester.TestService(testRunner, testRole, writer, shutdownChannel)
		} else {
			testResult, err = tester.TestApp(appContext, writer, buildMode, shutdownChannel)
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
	testCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run tests without mounting")
}
