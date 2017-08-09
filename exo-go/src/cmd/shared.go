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

func getAwsConfig() (types.AwsConfig, error) {
	appConfig, err := getAppConfig()
	if err != nil {
		return types.AwsConfig{}, err
	}
	return types.AwsConfig{
		Region:               "us-west-2", //TODO: read from app.yml
		SecretsBucket:        fmt.Sprintf("%s-terraform-secrets", appConfig.Name),
		TerraformStateBucket: fmt.Sprintf("%s-terraform", appConfig.Name),
		TerraformLockTable:   "TerraformLocks",
	}, nil
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
