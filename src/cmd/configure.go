package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	prompt "github.com/kofalt/go-prompt"
	"github.com/spf13/cobra"
)

var configureProfileFlag string

var configureCmd = &cobra.Command{
	Use:   "configure [remote-environment-id]",
	Short: "Configures secrets for an Exosphere application deployed to the cloud",
}

var configureReadCmd = &cobra.Command{
	Use:   "read [remote-environment-id]",
	Short: "Reads and prints secrets from remote secrets store",
	Long:  "Reads and prints secrets from remote secrets store",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Reading secrets store...\n\n")
		deployConfig := getBaseDeployConfig(args[0], configureProfileFlag)
		secrets, err := aws.ReadSecrets(deployConfig.GetAwsOptions())
		if err != nil {
			log.Fatalf("Cannot read secrets: %s", err)
		}
		prettyPrintSecrets(secrets)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

var configureCreateCmd = &cobra.Command{
	Use:   "create [remote-environment-id]",
	Short: "Creates a secret key entries in remote secrets store",
	Long:  "Creates a secret key entries in remote secrets store. Cannot conflict with existing keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("We are about to add secrets to the secret store!\n\n")
		deployConfig := getBaseDeployConfig(args[0], configureProfileFlag)
		existingSecrets := getSecrets(deployConfig)
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
			fmt.Print("\nYou are creating these secrets:\n\n")
			prettyPrintSecrets(newSecrets)

			if ok := prompt.Confirm("Do you want to continue? (y/n)"); ok {
				err := aws.MergeAndWriteSecrets(existingSecrets, newSecrets, deployConfig.GetAwsOptions())
				if err != nil {
					log.Fatalf("Cannot create secrets: %s", err)
				}
			} else {
				fmt.Println("Secret creation abandoned.")
			}
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

var configureUpdateCmd = &cobra.Command{
	Use:   "update [remote-environment-id]",
	Short: "Updates secret key entries in remote secrets store",
	Long:  "Updates secret key entries in remote secret store. Keys should already exist.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("We are about update keys in the remote store!\n\n")
		deployConfig := getBaseDeployConfig(args[0], configureProfileFlag)
		existingSecrets := getSecrets(deployConfig)
		existingSecretKeys := existingSecrets.Keys()
		newSecrets := map[string]string{}
		ok := true
		for ok {
			i := prompt.Choose("Select secret keys to update", existingSecretKeys)
			value := prompt.StringRequired(fmt.Sprintf("Secret value for %s", existingSecretKeys[i]))
			newSecrets[existingSecretKeys[i]] = value
			existingSecretKeys = append(existingSecretKeys[:i], existingSecretKeys[i+1:]...)
			ok = prompt.Confirm("Do you have more keys to update? (y/n)")
		}

		if len(newSecrets) > 0 {
			fmt.Print("\nYou are updating these secrets:\n\n")
			prettyPrintSecrets(newSecrets)

			if ok := prompt.Confirm("Do you want to continue? (y/n)"); ok {
				err := aws.MergeAndWriteSecrets(existingSecrets, newSecrets, deployConfig.GetAwsOptions())
				if err != nil {
					log.Fatalf("Cannot update secrets: %s", err)
				}
			} else {
				fmt.Println("Secret update abandoned.")
			}
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

var configureDeleteCmd = &cobra.Command{
	Use:   "delete [remote-environment-id]",
	Short: "Deletes secrets from the remote secrets store",
	Long:  "Deletes secrets from the remote secrets store. Ignores any keys passed in that don't exist on the remote store.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("We are about to delete secrets from the secret store...\n\n")
		deployConfig := getBaseDeployConfig(args[0], configureProfileFlag)
		existingSecrets := getSecrets(deployConfig)
		existingSecretKeys := existingSecrets.Keys()
		secretKeys := []string{}
		ok := true
		for ok {
			i := prompt.Choose("Select secret keys to delete", existingSecretKeys)
			secretKeys = append(secretKeys, existingSecretKeys[i])
			existingSecretKeys = append(existingSecretKeys[:i], existingSecretKeys[i+1:]...)
			ok = prompt.Confirm("Do you have more keys to delete? (y/n)")
		}

		if len(secretKeys) > 0 {
			fmt.Print("\nYou are deleting these secrets:\n\n")
			fmt.Printf("%s\n\n", strings.Join(secretKeys, ", "))

			if ok := prompt.Confirm("Do you want to continue? (y/n)"); ok {
				err := aws.DeleteSecrets(existingSecrets, secretKeys, deployConfig.GetAwsOptions())
				if err != nil {
					log.Fatalf("Cannot delete secrets: %s", err)
				}
			} else {
				fmt.Println("Secret deletion abandoned.")
			}
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateArgCount(args, 1)
	},
}

func init() {
	configureCmd.AddCommand(configureReadCmd)
	configureCmd.AddCommand(configureCreateCmd)
	configureCmd.AddCommand(configureUpdateCmd)
	configureCmd.AddCommand(configureDeleteCmd)
	RootCmd.AddCommand(configureCmd)
	configureCmd.PersistentFlags().StringVarP(&configureProfileFlag, "profile", "p", "default", "AWS profile to use")
}
