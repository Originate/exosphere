package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/aws_helpers"
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
		fmt.Print("We are about to configure the secrets store!\n\n")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
		err = awsHelper.CreateSecretsStore(secretsBucket, awsRegion)
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
		fmt.Println("Secrets store configured!")
	},
}

var configureReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads and prints secrets from remote secrets store",
	Long:  "Reads and prints secrets from remote secrets store",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("Reading secrets store...\n\n")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
		secrets, err := awsHelper.ReadSecrets(secretsBucket, awsRegion)
		if err != nil {
			log.Fatalf("Cannot read secrets: %s", err)
		}
		fmt.Println(secrets)
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

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
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
			log.Fatalf("Could not marshal secrets map: %s", err)
		}
		fmt.Print("You are creating these secrets:\n\n")
		fmt.Printf("%s\n\n", string(secretsPretty))

		if ok := prompt.Confirm("Do you want to continue?"); ok {
			err = awsHelper.CreateSecrets(secrets, secretsBucket, awsRegion)
			if err != nil {
				log.Fatalf("Cannot create secrets: %s", err)
			}
		} else {
			fmt.Println("Secret creation abandoned.")
		}
	},
}

func init() {
	configureCmd.AddCommand(configureReadCmd)
	configureCmd.AddCommand(configureCreateCmd)
	RootCmd.AddCommand(configureCmd)
}
