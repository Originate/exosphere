package cmd

import (
	"log"
	"os"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/types/deploy"
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
			deployConfig := deploy.Config{
				AppContext:          userContext.AppContext,
				RemoteEnvironmentID: remoteEnvironmentID,
				Writer:              os.Stdout,
			}
			err = application.GenerateTerraformFiles(deployConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var generateTerraformVarFileCmd = &cobra.Command{
	Use:   "terraform-var-file [remote-environment-id]",
	Short: "Generates terraform tfvar files",
	Long:  "Generates terraform tfvar files",
	Run: func(cmd *cobra.Command, args []string) {
		deployConfig := getBaseDeployConfig(args[0], deployProfileFlag)
		secrets, err := aws.ReadSecrets(deployConfig.GetAwsOptions())
		if err != nil {
			log.Fatalf("Could not read aws secrets: %s", err)
		}
		err = application.GenerateTerraformVarFiles(deployConfig, secrets)
		if err != nil {
			log.Fatal(err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

func init() {
	generateDockerComposeCmd.PersistentFlags().BoolVarP(&checkFlag, "check", "", false, "Runs check to see if docker-compose are up-to-date")
	generateTerraformVarFileCmd.PersistentFlags().StringVarP(&deployProfileFlag, "profile", "p", "default", "AWS profile to use")
	generateCmd.AddCommand(generateDockerComposeCmd)
	generateCmd.AddCommand(generateTerraformCmd)
	generateCmd.AddCommand(generateTerraformVarFileCmd)
	RootCmd.AddCommand(generateCmd)
}
