package cmd

import (
	"fmt"
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
		fmt.Print("Generating docker-compose and terraform files...\n\n")
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
		fmt.Println("docker-compose and terraform files generated")
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}
