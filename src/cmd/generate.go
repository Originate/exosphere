package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/aws"
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
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		if len(userContext.AppContext.Config.Remote.Environments) == 0 {
			log.Fatal("No remote environments. Add one in your `application.yml` at `remote.environments.<id>`")
		}
		for remoteEnvironmentID := range userContext.AppContext.Config.Remote.Environments {
			deployConfig := getBaseDeployConfig(userContext.AppContext, remoteEnvironmentID)
			deployConfig.Writer = os.Stdout
			err = application.GenerateTerraformFiles(deployConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var generateTerraformTfVarsCmd = &cobra.Command{
	Use:   "terraform-var-files",
	Short: "Generates terraform tfvar files",
	Long:  "Generates terraform tfvar files",
	Run: func(cmd *cobra.Command, args []string) {
		userContext, err := GetUserContext()
		if err != nil {
			log.Fatal(err)
		}
		if len(userContext.AppContext.Config.Remote.Environments) == 0 {
			log.Fatal("No remote environments. Add one in your `application.yml` at `remote.environments.<id>`")
		}
		for remoteEnvironmentID := range userContext.AppContext.Config.Remote.Environments {
			deployConfig := getBaseDeployConfig(userContext.AppContext, remoteEnvironmentID)
			deployConfig.Writer = os.Stdout
			secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
			if err != nil {
				log.Fatalf("Could not read aws secrets: %s", err)
			}
			err = application.GenerateTerraformVarFile(deployConfig, secrets)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	generateDockerComposeCmd.PersistentFlags().BoolVarP(&checkFlag, "check", "", false, "Runs check to see if docker-compose are up-to-date")
	generateTerraformTfVarsCmd.PersistentFlags().StringVarP(&awsProfileFlag, "profile", "p", "default", "AWS profile to use")
	generateCmd.AddCommand(generateDockerComposeCmd)
	generateCmd.AddCommand(generateTerraformCmd)
	generateCmd.AddCommand(generateTerraformTfVarsCmd)
	RootCmd.AddCommand(generateCmd)
}
