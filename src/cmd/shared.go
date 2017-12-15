package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
)

func getAwsConfig(appConfig types.AppConfig, profile string) types.AwsConfig {
	return types.AwsConfig{
		Region:               appConfig.Remote.Region,
		AccountID:            appConfig.Remote.AccountID,
		SslCertificateArn:    appConfig.Remote.SslCertificateArn,
		Profile:              profile,
		SecretsBucket:        fmt.Sprintf("%s-%s-terraform-secrets", appConfig.Remote.AccountID, appConfig.Name),
		TerraformStateBucket: fmt.Sprintf("%s-%s-terraform", appConfig.Remote.AccountID, appConfig.Name),
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

func getBaseDeployConfig(appContext *context.AppContext) deploy.Config {
	awsConfig := getAwsConfig(appContext.Config, deployProfileFlag)
	return deploy.Config{
		AppContext: appContext,
		AwsConfig:  awsConfig,
	}
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}
