package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/aws_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/segmentio/go-prompt"
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

		secretsBucket, awsRegion := getBucketConfig()
		err := awsHelper.CreateSecretsStore(secretsBucket, awsRegion)
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
		fmt.Println("We are about to add secrets to the secret store!")

		secretsBucket, awsRegion := getBucketConfig()
		secrets := map[string]string{}

		secretName := prompt.String("Secret name")
		if secretName != "" {
			secretValue := prompt.StringRequired("Secret value")
			secrets[secretName] = secretValue
		}

		for secretName != "" {
			secretName = prompt.String("Secret name (leave blank to finish prompting)")
			if secretName != "" {
				secretValue := prompt.StringRequired("Secret value")
				secrets[secretName] = secretValue
			}
		}

		secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
		if err != nil {
			log.Fatalf("Could not marshal secrets map", err)
		}
		fmt.Println("You are creating these secrets:\n")
		fmt.Printf("%s\n\n", string(secretsPretty))

		if ok := prompt.Confirm("Do you want to continue?"); ok {
			awsHelper.CreateSecrets(secrets, secretsBucket, awsRegion)
		} else {
			fmt.Println("Secret creation abandoned.")
		}
	},
}

func getAppConfig() types.AppConfig {
	appDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfig, err := appConfigHelpers.GetAppConfig(appDir)
	if err != nil {
		log.Fatalf("Cannot read application configuration: %s", err)
	}
	return appConfig
}

func getBucketConfig() (string, string) {
	appConfig := getAppConfig()
	secretsBucket := fmt.Sprintf("%s-terraform-secrets", appConfig.Name)
	awsRegion := "us-west-2" //TODO: read from app.yml
	return secretsBucket, awsRegion
}

func init() {
	configureCmd.AddCommand(configureCreateCmd)
	RootCmd.AddCommand(configureCmd)
}
