package cmd

import (
	"fmt"
	"log"
	"os"

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

		context, err := types.GetContext()
		if err != nil {
			log.Fatal(err)
		}
		context.AppContext.DockerComposeProjectName = fmt.Sprintf("%stests", composebuilder.GetDockerComposeProjectName(context.AppContext.Location))
		serviceRoles := context.AppContext.Config.GetSortedServiceRoles()
		dependencyNames := context.AppContext.Config.GetDevelopmentDependencyNames()
		roles := append(serviceRoles, dependencyNames...)
		roles = append(roles, "exo-test")
		logger := util.NewLogger(roles, []string{}, "exo-test", os.Stdout)
		buildMode := composebuilder.BuildModeLocalDevelopment
		if noMountFlag {
			buildMode = composebuilder.BuildModeLocalDevelopmentNoMount
		}

		var testsPassed bool
		var testsInterrupted bool
		if context.HasServiceContext {
			testsPassed, testsInterrupted, err = application.TestService(context.ServiceContext, logger, buildMode)
		} else {
			testsPassed, testsInterrupted, err = application.TestApp(context.AppContext, logger, buildMode)
		}
		if err != nil {
			panic(err)
		}
		if !testsPassed && !testsInterrupted {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run tests without mounting")
}
