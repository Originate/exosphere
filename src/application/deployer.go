package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	fmt.Printf("Setting up AWS account...\n\n")
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	fmt.Printf("Building %s %s...\n\n", deployConfig.AppConfig.Name, deployConfig.AppConfig.Version)
	initializer, err := NewInitializer(deployConfig.AppConfig, deployConfig.Logger, "exo-deploy", deployConfig.AppDir, deployConfig.HomeDir)
	if err != nil {
		return err
	}
	err = initializer.Initialize()
	if err != nil {
		return err
	}

	fmt.Printf("\n\nPushing Docker images to ECR...\n\n")
	dockerComposePath := filepath.Join(deployConfig.AppDir, "tmp", "docker-compose.yml")
	imagesMap, err := aws.PushImages(deployConfig, dockerComposePath)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nGenerating Terraform files...\n\n")
	err = terraform.GenerateFile(deployConfig, imagesMap)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving remote state...\n\n")
	err = terraform.RunInit(deployConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving secrets...\n\n")
	err = writeSecretsFile(deployConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nPlanning deployment...\n\n")
	return terraform.RunPlan(deployConfig)
}

func writeSecretsFile(deployConfig types.DeployConfig) error {
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	if err != nil {
		return fmt.Errorf("Cannot read secrets: %s", err)
	}
	var filePerm os.FileMode = 0644 //standard Unix file permission: rw-rw-rw-
	return ioutil.WriteFile(deployConfig.SecretsPath, []byte(secrets), filePerm)
}

// RemoveSecretsFile removes the secrets file from the user's machine
func RemoveSecretsFile(secretsPath string) error {
	err := os.Remove(secretsPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error removing secrets file: %s. Manual removal recommended", secretsPath))
	}
	return nil
}
