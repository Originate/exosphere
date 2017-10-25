package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var noMountFlag bool

func printHeader(text string) {
	separator := strings.Repeat("*", 80)
	fmt.Printf("\n%s\n%s\n%s\n", separator, text, separator)
}

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

func getAwsConfig(profile string) (types.AwsConfig, error) {
	appConfig, err := getAppConfig()
	if err != nil {
		return types.AwsConfig{}, err
	}
	return types.AwsConfig{
		Region:               appConfig.Production.Region,
		AccountID:            appConfig.Production.AccountID,
		SslCertificateArn:    appConfig.Production.SslCertificateArn,
		Profile:              profile,
		SecretsBucket:        fmt.Sprintf("%s-%s-terraform-secrets", appConfig.Production.AccountID, appConfig.Name),
		TerraformStateBucket: fmt.Sprintf("%s-%s-terraform", appConfig.Production.AccountID, appConfig.Name),
		TerraformLockTable:   "TerraformLocks",
	}, nil
}

func getSecrets(awsConfig types.AwsConfig) types.Secrets {
	secrets, err := aws.ReadSecrets(awsConfig)
	if err != nil {
		log.Fatalf("Cannot read secrets: %s", err)
	}
	fmt.Print("Existing secrets:\n\n")
	prettyPrintSecrets(secrets)
	return secrets
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}

// returns the name for a test project
func getTestDockerComposeProjectName(appDir string) string {
	return fmt.Sprintf("%stests", composebuilder.GetDockerComposeProjectName(appDir))
}
