package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
		fmt.Println("We are about to add secrets to the secret store!")

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
		secretName := prompt.String("Secret name")
		if secretName != "" {
			if _, hasKey := existingSecrets[secretName]; hasKey {
				fmt.Printf("Secret for '%s' already exists. Use 'exo configure update' to update.\n\n", secretName)
			} else {
				secretValue := prompt.StringRequired("Secret value")
				newSecrets[secretName] = secretValue
			}
		}

		for secretName != "" {
			secretName = prompt.String("Secret name (leave blank to finish prompting)")
			if secretName != "" {
				if _, hasKey := existingSecrets[secretName]; hasKey {
					fmt.Printf("Secret for '%s' already exists. Use 'exo configure update' to update.\n\n", secretName)
				} else {
					secretValue := prompt.StringRequired("Secret value")
					newSecrets[secretName] = secretValue
				}
			}
		}

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
	},
}

var configureUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates secret key entries in remote secrets store",
	Long:  "Updates secret key entries in remote secret store. Keys should already exist.",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Println("We are about update keys in the remote store!")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot update secrets: %s", err)
		}

		secretsString, err := aws.ReadSecrets(secretsBucket, awsRegion)
		if err != nil {
			log.Fatalf("Cannot read secrets: %s", err)
		}
		existingSecrets := types.NewSecrets(secretsString)

		newSecrets := map[string]string{}
		secretName := prompt.String("Secret name")
		if secretName != "" {
			if _, hasKey := existingSecrets[secretName]; !hasKey {
				fmt.Printf("Secret for '%s' does not exists. Use 'exo configure create' to create it.\n\n", secretName)
			} else {
				secretValue := prompt.StringRequired("Secret value")
				newSecrets[secretName] = secretValue
			}
		}
		for secretName != "" {
			secretName = prompt.String("Secret name (leave blank to finish prompting)")
			if secretName != "" {
				if _, hasKey := existingSecrets[secretName]; !hasKey {
					fmt.Printf("Secret for '%s' does not exists. Use 'exo configure create' to create it.\n\n", secretName)
				} else {
					secretValue := prompt.StringRequired("Secret value")
					newSecrets[secretName] = secretValue
				}
			}
		}

		secretsPretty, err := json.MarshalIndent(newSecrets, "", "  ")
		if err != nil {
			log.Fatalf("Could not marshal secrets map: %s", err)
		}
		fmt.Print("You are updating these secrets:\n\n")
		fmt.Printf("%s\n\n", string(secretsPretty))

		if ok := prompt.Confirm("Do you want to continue?"); ok {
			err = aws.MergeAndWriteSecrets(existingSecrets, newSecrets, secretsBucket, awsRegion)
			if err != nil {
				log.Fatalf("Cannot update secrets: %s", err)
			}
		} else {
			fmt.Println("Secret update abandoned.")
		}
	},
}

var configureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes secrets from the remote secrets store",
	Long:  "Deletes secrets from the remote secrets store. Ignores any keys passed in that don't exist on the remote store.",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		fmt.Print("We are about to delete secrets from the secret store...\n\n")

		secretsBucket, awsRegion, err := getSecretsConfig()
		if err != nil {
			log.Fatalf("Cannot delete secrets: %s", err)
		}

		secretKeys := []string{}
		secretName := prompt.String("Secret name")
		if secretName != "" {
			secretKeys = append(secretKeys, secretName)
		}
		for secretName != "" {
			secretName = prompt.String("Secret name (leave blank to finish prompting)")
			if secretName != "" {
				secretKeys = append(secretKeys, secretName)
			}
		}

		fmt.Print("You are deleting these secrets:\n\n")
		fmt.Printf("%s\n\n", strings.Join(secretKeys, ", "))

		if ok := prompt.Confirm("Do you want to continue?"); ok {
			err := aws.DeleteSecrets(secretKeys, secretsBucket, awsRegion)
			if err != nil {
				log.Fatalf("Cannot delete secrets: %s", err)
			}
		} else {
			fmt.Println("Secret deletion abandoned.")
		}
	},
}

func init() {
	configureCmd.AddCommand(configureReadCmd)
	configureCmd.AddCommand(configureCreateCmd)
	configureCmd.AddCommand(configureUpdateCmd)
	configureCmd.AddCommand(configureDeleteCmd)
	RootCmd.AddCommand(configureCmd)
}
