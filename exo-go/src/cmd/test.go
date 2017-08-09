package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-go/src/application"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
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
		if _, err := os.Stat("application.yml"); err == nil {
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
			serviceNames := appConfig.GetServiceNames()
			dependencyNames := appConfig.GetDependencyNames()
			roles := append(serviceNames, dependencyNames...)
			roles = append(roles, "exo-test")
			logger := application.NewLogger(roles, []string{}, os.Stdout)
			tester, err := application.NewTester(appConfig, logger, appDir, homeDir)
			if err != nil {
				panic(err)
			}
			if err := tester.RunAppTests(); err != nil {
				panic(err)
			}
		} else if _, err := os.Stat("service.yml"); err == nil {
			fmt.Println("This is trueeeeee")
		} else {
			fmt.Println("Not an application or service directory, exiting...")
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
