package cmd

import (
	"log"

	"github.com/Originate/exosphere/src/application"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates docker-compose and terraform files",
	Long:  "Generates docker-compose and terraform files",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		context, err := GetContext()
		if err != nil {
			log.Fatal(err)
		}
		err = application.GenerateComposeFiles(context.AppContext)
		if err != nil {
			log.Fatal(err)
		}
		deployConfig := getBaseDeployConfig(context.AppContext)
		err = application.GenerateTerraformFiles(deployConfig)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
