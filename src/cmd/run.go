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
		dockerComposeFileName := types.LocalDevelopmentComposeFileName
		if productionFlag {
			dockerComposeFileName = types.LocalProductionComposeFileName
		}
		err = runner.Run(runner.RunOptions{
			AppContext:               userContext.AppContext,
			DockerComposeFileName:    dockerComposeFileName,
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
