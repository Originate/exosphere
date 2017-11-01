package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/util"
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
		context, err := GetContext()
		if err != nil {
			log.Fatal(err)
		}
		homeDir, err := util.GetHomeDirectory()
		if err != nil {
			panic(err)
		}
		dockerComposeProjectName := composebuilder.GetDockerComposeProjectName(context.AppContext.Config.Name)
		writer := os.Stdout
		buildMode := composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Mount:       true,
			Environment: composebuilder.BuildModeEnvironmentDevelopment,
		}
		if productionFlag {
			buildMode.Environment = composebuilder.BuildModeEnvironmentProduction
		} else if noMountFlag {
			buildMode.Mount = false
		}
		runner, err := application.NewRunner(context.AppContext, writer, homeDir, dockerComposeProjectName, buildMode)
		if err != nil {
			panic(err)
		}
		if err := runner.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&noMountFlag, "no-mount", "", false, "Run without mounting")
	runCmd.PersistentFlags().BoolVarP(&productionFlag, "production", "", false, "Run in production mode")
}
