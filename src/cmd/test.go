package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
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
		var appDir, serviceName string
		var testsPassed bool
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if _, err = os.Stat("service.yml"); err == nil {
			serviceConfig := types.ServiceConfig{}
			var yamlFile []byte
			yamlFile, err = ioutil.ReadFile(path.Join(currentDir, "service.yml"))
			if err != nil {
				panic(err)
			}
			if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
				panic(err)
			}
			serviceName = serviceConfig.Type
			appDir = path.Join(currentDir, "/..")
		} else if _, err = os.Stat("application.yml"); err == nil {
			appDir = currentDir
		} else {
			fmt.Println("Not an application or service directory, exiting...")
			os.Exit(1)
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
		dependencyNames := appConfig.GetDevelopmentDependencyNames()
		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-test")
		logger := application.NewLogger(roles, []string{}, os.Stdout)
		logChannel := logger.GetLogChannel("exo-test")
		dockerComposeProjectName := getTestDockerComposeProjectName(appDir)
		buildMode := composebuilder.BuildModeLocalDevelopment
		if noMountFlag {
			buildMode = composebuilder.BuildModeLocalDevelopmentNoMount
		}
		tester, err := application.NewTester(appConfig, logChannel, appDir, homeDir, dockerComposeProjectName, buildMode)
		if err != nil {
			panic(err)
		}
		if serviceName != "" {
			if testsPassed, err = tester.RunServiceTest(serviceName); err != nil {
				panic(err)
			}
		} else {
			if testsPassed, err = tester.RunAppTests(); err != nil {
				panic(err)
			}
		}
		close(logChannel)
		logger.WaitForChannelsToClose()
		if !testsPassed {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run tests without mounting")
}
