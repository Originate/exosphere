package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
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

func getAwsConfig(appConfig types.AppConfig, remoteID string, profile string) types.AwsConfig {
	remoteConfig := appConfig.Remote[remoteID]
	return types.AwsConfig{
		Region:               remoteConfig.Region,
		AccountID:            remoteConfig.AccountID,
		SslCertificateArn:    remoteConfig.SslCertificateArn,
		Profile:              profile,
		SecretsBucket:        fmt.Sprintf("%s-%s-terraform-secrets", remoteConfig.AccountID, appConfig.Name),
		TerraformStateBucket: fmt.Sprintf("%s-%s-terraform", remoteConfig.AccountID, appConfig.Name),
		TerraformLockTable:   "TerraformLocks",
	}
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

func getBaseDeployConfig(appContext *context.AppContext, remoteID string) deploy.Config {
	awsConfig := getAwsConfig(appContext.Config, remoteID, deployProfileFlag)
	return deploy.Config{
		AppContext: appContext,
		AwsConfig:  awsConfig,
		RemoteID:   remoteID,
	}
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}

func validateArgCount(args []string, number int) error {
	if len(args) != number {
		return errors.New("Wrong number of arguments")
	}
	return nil
}

func validateRemoteID(userContext *context.UserContext, remoteID string) error {
	if _, ok := userContext.AppContext.Config.Remote[remoteID]; !ok {
		validRemoteIDs := []string{}
		for remoteID := range userContext.AppContext.Config.Remote {
			validRemoteIDs = append(validRemoteIDs, remoteID)
		}
		return fmt.Errorf("Invalid remote id: %s. Valid remote ids: %s", remoteID, strings.Join(validRemoteIDs, ","))
	}
	return nil
}
