package cmd

import (
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
		serviceNames := context.AppContext.Config.GetSortedServiceNames()
		dependencyNames := context.AppContext.Config.GetDevelopmentDependencyNames()
		roles := append(serviceNames, dependencyNames...)
		roles = append(roles, "exo-test")
		logger := util.NewLogger(roles, []string{}, "exo-test", os.Stdout)
		buildMode := composebuilder.BuildModeLocalDevelopment
		if noMountFlag {
			buildMode = composebuilder.BuildModeLocalDevelopmentNoMount
		}

		var testsPassed bool
		if context.HasServiceContext {
			testsPassed, err = application.TestService(context.ServiceContext, logger, buildMode)
		} else {
			testsPassed, err = application.TestApp(context.AppContext, logger, buildMode)
		}
		if err != nil {
			panic(err)
		}
		if !testsPassed {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run tests without mounting")
}
