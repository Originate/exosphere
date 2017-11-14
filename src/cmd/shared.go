package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/spf13/cobra"
)

var noMountFlag bool

func printHelpIfNecessary(cmd *cobra.Command, args []string) bool {
	if len(args) == 1 && args[0] == "help" {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func getAwsConfig(appConfig types.AppConfig, profile string) types.AwsConfig {
	return types.AwsConfig{
		Region:               appConfig.Production.Region,
		AccountID:            appConfig.Production.AccountID,
		SslCertificateArn:    appConfig.Production.SslCertificateArn,
		Profile:              profile,
		SecretsBucket:        fmt.Sprintf("%s-%s-terraform-secrets", appConfig.Production.AccountID, appConfig.Name),
		TerraformStateBucket: fmt.Sprintf("%s-%s-terraform", appConfig.Production.AccountID, appConfig.Name),
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

func getBaseDeployConfig(appContext types.AppContext) types.DeployConfig {
	serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
	if err != nil {
		log.Fatalf("Failed to read service configurations: %s", err)
	}
	awsConfig := getAwsConfig(appContext.Config, deployProfileFlag)
	terraformDir := filepath.Join(appContext.Location, "terraform")
	return types.DeployConfig{
		AppContext:               appContext,
		ServiceConfigs:           serviceConfigs,
		DockerComposeProjectName: composebuilder.GetDockerComposeProjectName(appContext.Config.Name),
		TerraformDir:             terraformDir,
		SecretsPath:              filepath.Join(terraformDir, "secrets.tfvars"),
		AwsConfig:                awsConfig,

		// git commit hash of the Terraform modules in Originate/exosphere we are using
		TerraformModulesRef: "1bf0375f",
	}
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}
