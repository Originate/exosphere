package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/aws_helpers"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configures secrets for an Exosphere application deployed to the cloud",
	Long:  "Configures secrets for an Exosphere application deployed to the cloud. Creates a remote secret store",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}

		appDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		appConfig, err := appConfigHelpers.GetAppConfig(appDir)
		if err != nil {
			log.Fatalf("Cannot read application configuration: %s", err)
		}

		secretsBucket := fmt.Sprintf("%s-terraform-secrets", appConfig.Name)
		region := "us-west-2" //TODO: prompt user for this and consolidate with exo-deploy

		err = awsHelper.CreateSecretsStore(secretsBucket, region)
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
	},
}

var configureCreateCmd = &cobra.Command{
	Use:   "create [secrets]",
	Short: "Creates a secret key entry in remote secrets store",
	Long:  "Creates a secret key entry in the remote secrets store. Should be in the form 'secret_key=secret_value'",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		awsHelper.CreateSecretEntry(args)
	},
}

func init() {
	configureCmd.AddCommand(configureCreateCmd)
	RootCmd.AddCommand(configureCmd)
}
