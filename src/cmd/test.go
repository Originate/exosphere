package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/tester"
	"github.com/Originate/exosphere/src/types"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs tests for the application",
	Long:  "Runs tests for the application",
	Run: func(cmd *cobra.Command, args []string) {
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		err = application.GenerateComposeFiles(userContext.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		writer := os.Stdout
		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)
		var testResult types.TestResult
		if userContext.HasServiceContext {
			testResult, err = tester.TestService(userContext.ServiceContext, writer, shutdownChannel)
		} else {
			testResult, err = tester.TestApp(userContext.AppContext, writer, shutdownChannel)
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
