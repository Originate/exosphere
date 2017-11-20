package cmd

import (
	"log"

	"github.com/Originate/exosphere/src/application"
	"github.com/spf13/cobra"
)

var generateCheckFlag bool

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
		deployConfig := getBaseDeployConfig(context.AppContext)
		if generateCheckFlag {
			err = application.CheckGeneratedFiles(context.AppContext, deployConfig)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = application.GenerateComposeFiles(context.AppContext)
			if err != nil {
				log.Fatal(err)
			}
			err = application.GenerateTerraformFiles(deployConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().BoolVarP(&generateCheckFlag, "check", "", false, "Runs check to see if docker-compose and terraform files are up-to-date")
}
