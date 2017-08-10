package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
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
		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		serviceName := ""
		if _, err = os.Stat("service.yml"); err == nil {
			serviceConfig := types.ServiceConfig{}
			yamlFile, err := ioutil.ReadFile(path.Join(appDir, "service.yml"))
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(yamlFile, &serviceConfig)
			serviceName = serviceConfig.Type
			appDir = path.Join(appDir, "/..")
		} else if _, err = os.Stat("application.yml"); os.IsNotExist(err) {
			fmt.Println("Not an application or service directory, exiting...")
			return
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
		if serviceName != "" {
			if err := tester.RunServiceTest(serviceName); err != nil {
				panic(err)
			}
		} else {
			if err := tester.RunAppTests(); err != nil {
				panic(err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
