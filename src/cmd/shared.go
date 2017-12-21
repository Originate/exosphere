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
)

func getAwsConfig(appConfig types.AppConfig, remoteEnvironmentID string, profile string) types.AwsConfig {
	appRemoteEnv := appConfig.Remote.Environments[remoteEnvironmentID]
	return types.AwsConfig{
		Region:               appRemoteEnv.Region,
		AccountID:            appRemoteEnv.AccountID,
		SslCertificateArn:    appRemoteEnv.SslCertificateArn,
		Profile:              profile,
		SecretsBucket:        fmt.Sprintf("%s-%s-%s-terraform-secrets", appRemoteEnv.AccountID, appConfig.Name, remoteEnvironmentID),
		TerraformStateBucket: fmt.Sprintf("%s-%s-%s-terraform", appRemoteEnv.AccountID, appConfig.Name, remoteEnvironmentID),
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

func getBaseDeployConfig(appContext *context.AppContext, remoteEnvironmentID string) deploy.Config {
	awsConfig := getAwsConfig(appContext.Config, remoteEnvironmentID, deployProfileFlag)
	return deploy.Config{
		AppContext:          appContext,
		AwsConfig:           awsConfig,
		RemoteEnvironmentID: remoteEnvironmentID,
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

func validateRemoteEnvironmentID(userContext *context.UserContext, remoteEnvironmentID string) error {
	if _, ok := userContext.AppContext.Config.Remote.Environments[remoteEnvironmentID]; !ok {
		validIDs := []string{}
		for validID := range userContext.AppContext.Config.Remote.Environments {
			validIDs = append(validIDs, validID)
		}
		return fmt.Errorf("Invalid remote environment id: %s. Valid remote environment ids: %s", remoteEnvironmentID, strings.Join(validIDs, ","))
	}
	return nil
}
