package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composebuilder"
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
		dockerComposeProjectName := composebuilder.GetDockerComposeProjectName(context.AppContext.Config.Name)
		writer := os.Stdout
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentDevelopment,
		}
		if productionFlag {
			buildMode.Environment = types.BuildModeEnvironmentProduction
		}
		runner, err := application.NewRunner(context.AppContext, writer, dockerComposeProjectName, buildMode)
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
	runCmd.PersistentFlags().BoolVarP(&productionFlag, "production", "", false, "Run in production mode")
}
