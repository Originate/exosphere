package cmd

import (
	"log"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/spf13/cobra"
)

var checkFlag bool

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates docker-compose and terraform files",
	Long:  "Generates docker-compose and terraform files",
}

var generateDockerComposeCmd = &cobra.Command{
	Use:   "docker-compose",
	Short: "Generates docker-compose files",
	Long:  "Generates docker-compose files",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		if checkFlag {
			err = application.CheckGeneratedDockerComposeFiles(userContext.AppContext)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = application.GenerateComposeFiles(userContext.AppContext)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var generateTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Generates terraform files",
	Long:  "Generates terraform files",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		deployConfig := getBaseDeployConfig(userContext.AppContext)
		err = terraform.GenerateFile(deployConfig)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	generateDockerComposeCmd.PersistentFlags().BoolVarP(&checkFlag, "check", "", false, "Runs check to see if docker-compose are up-to-date")
	generateCmd.AddCommand(generateDockerComposeCmd)
	generateCmd.AddCommand(generateTerraformCmd)
	RootCmd.AddCommand(generateCmd)
}
