package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/application/runner"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/spf13/cobra"
)

var productionFlag bool

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an Exosphere application",
	Long:  "Runs an Exosphere application",
	Run: func(cmd *cobra.Command, args []string) {
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		err = application.GenerateComposeFiles(userContext.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		buildMode := types.BuildMode{
			Type:        types.BuildModeTypeLocal,
			Mount:       true,
			Environment: types.BuildModeEnvironmentDevelopment,
		}
		if productionFlag {
			buildMode.Environment = types.BuildModeEnvironmentProduction
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
