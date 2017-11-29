package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application/runner"
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
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		buildMode := composebuilder.BuildMode{
			Type:        composebuilder.BuildModeTypeLocal,
			Mount:       true,
			Environment: composebuilder.BuildModeEnvironmentDevelopment,
		}
		if productionFlag {
			buildMode.Environment = composebuilder.BuildModeEnvironmentProduction
		}
		err = runner.Run(runner.RunOptions{
			AppContext:               userContext.AppContext,
			BuildMode:                buildMode,
			DockerComposeProjectName: composebuilder.GetDockerComposeProjectName(userContext.AppContext.Config.Name),
			Writer: os.Stdout,
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVarP(&productionFlag, "production", "", false, "Run in production mode")
}
