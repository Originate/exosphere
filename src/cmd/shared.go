package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
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

func getDeployConfig() (types.DeployConfig, error) {
	appDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	homeDir, err := util.GetHomeDirectory()
	if err != nil {
		panic(err)
	}
	appConfig, err := types.NewAppConfig(appDir)
	if err != nil {
		return types.DeployConfig{}, fmt.Errorf("Cannot read application configuration: %s", err)
	}
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return types.DeployConfig{}, fmt.Errorf("Cannot read service configurations: %s", err)
	}
	awsConfig, err := getAwsConfig()
	if err != nil {
		return types.DeployConfig{}, fmt.Errorf("Cannot read AWS configuration: %s", err)
	}

	logger := application.NewLogger([]string{"exo-deploy"}, []string{}, os.Stdout)
	terraformDir := filepath.Join(appDir, "terraform")
	deployConfig := types.DeployConfig{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
		LogChannel:     logger.GetLogChannel("exo-deploy"),
		TerraformDir:   terraformDir,
		TerraformFile:  filepath.Join(terraformDir, "main.tf"),
		SecretsPath:    filepath.Join(terraformDir, "secrets.tfvars"),
		AwsConfig:      awsConfig,
	}

	return deployConfig, nil

}

func getSecrets(awsConfig types.AwsConfig) types.Secrets {
	secretsString, err := aws.ReadSecrets(awsConfig)
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
