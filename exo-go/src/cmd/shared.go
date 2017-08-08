package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func printHelpIfNecessary(cmd *cobra.Command, args []string) bool {
	if len(args) == 1 && args[0] == "help" {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func getAppConfig() (types.AppConfig, error) {
	appDir, err := os.Getwd()
	if err != nil {
		return types.AppConfig{}, errors.Wrap(err, "Cannot get working directory")
	}
	appConfig, err := types.NewAppConfig(appDir)
	if err != nil {
		return types.AppConfig{}, errors.Wrap(err, "Cannot read application configuration")
	}
	return appConfig, nil
}

func getSecretsConfig() (string, string, error) {
	appConfig, err := getAppConfig()
	if err != nil {
		return "", "", err
	}
	secretsBucket := fmt.Sprintf("%s-terraform-secrets", appConfig.Name)
	awsRegion := "us-west-2" //TODO: read from app.yml
	return secretsBucket, awsRegion, nil
}

func getSecrets(secretsBucket, awsRegion string) types.Secrets {
	secretsString, err := aws.ReadSecrets(secretsBucket, awsRegion)
	if err != nil {
		log.Fatalf("Cannot read secrets: %s", err)
	}
	existingSecrets := types.NewSecrets(secretsString)
	fmt.Print("Existing secrets:\n\n")
	prettyPrintSecrets(existingSecrets)
	return existingSecrets
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}
