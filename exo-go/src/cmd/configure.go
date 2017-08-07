package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Originate/exosphere/exo-go/src/aws"
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
		fmt.Print("We are about to configure the secrets store!\n\n")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot create secrets store: %s", err)
		}
		err = aws.CreateSecretsStore(secretsBucket, awsRegion)
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
			log.Fatalf("Cannot read secrets: %s", err)
		}
		secrets, err := aws.ReadSecrets(secretsBucket, awsRegion)
		if err != nil {
			log.Fatalf("Cannot read secrets: %s", err)
		}
		fmt.Println(secrets)
	},
}

var configureCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a secret key entries in remote secrets store",
	Long:  "Creates a secret key entries in remote secrets store. Cannot conflict with existing keys",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to add secrets to the secret store!\n\n")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot create secrets: %s", err)
		}

		secretsString, err := aws.ReadSecrets(secretsBucket, awsRegion)
		if err != nil {
			log.Fatalf("Cannot read secrets: %s", err)
		}
		existingSecrets := types.NewSecrets(secretsString)

		newSecrets := map[string]string{}
		for {
			secretName := prompt.String("Secret name (leave blank to finish prompting)")
			if secretName == "" {
				break
			}
			if _, hasKey := existingSecrets[secretName]; hasKey {
				fmt.Printf("Secret for '%s' already exists. Use 'exo configure update' to update.\n\n", secretName)
			} else {
				secretValue := prompt.StringRequired("Secret value")
				newSecrets[secretName] = secretValue
			}
		}

		if len(newSecrets) > 0 {
			secretsPretty, err := json.MarshalIndent(newSecrets, "", "  ")
			if err != nil {
				log.Fatalf("Could not marshal secrets map: %s", err)
			}
			fmt.Print("You are creating these secrets:\n\n")
			fmt.Printf("%s\n\n", string(secretsPretty))

			if ok := prompt.Confirm("Do you want to continue?"); ok {
				err = aws.MergeAndWriteSecrets(existingSecrets, newSecrets, secretsBucket, awsRegion)
				if err != nil {
					log.Fatalf("Cannot create secrets: %s", err)
				}
			} else {
				fmt.Println("Secret creation abandoned.")
			}
		}
	},
}

func init() {
	configureCmd.AddCommand(configureReadCmd)
	configureCmd.AddCommand(configureCreateCmd)
	RootCmd.AddCommand(configureCmd)
}
