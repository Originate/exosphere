package application

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	fmt.Printf("Setting up AWS account...\n\n")
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	fmt.Printf("Building %s %s...\n\n", d.AppConfig.Name, d.AppConfig.Version)
	initializer, err := NewInitializer(d.AppConfig, d.Logger, "exo-deploy", d.AppDir, d.HomeDir)
	if err != nil {
		return err
	}
	err = initializer.Initialize()
	if err != nil {
		return err
	}

	fmt.Printf("\n\nPushing Docker images to ECR...\n\n")
	dockerComposePath := filepath.Join(d.AppDir, "tmp", "docker-compose.yml")
	err = aws.PushImages(d.AppConfig, dockerComposePath, terraformConfig.Region)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nGenerating Terraform files...\n\n")
	err = terraform.GenerateFile(deployConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving remote state...\n\n")
	err = terraform.RunInit(deployConfig.TerraformDir, deployConfig.Logger)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving secrets...\n\n")
	err = writeSecretsFile(deployConfig.SecretsPath, deployConfig.AwsConfig.SecretsBucket, deployConfig.AwsConfig.Region)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nPlanning deployment...\n\n")
	return terraform.RunPlan(deployConfig.TerraformDir, deployConfig.SecretsPath, deployConfig.Logger)
}

func writeSecretsFile(secretsPath, remoteBucket, region string) error {
	secrets, err := aws.ReadSecrets(remoteBucket, region)
	if err != nil {
		return fmt.Errorf("Cannot read secrets: %s", err)
	}
	var filePerm os.FileMode = 0644 //standard Unix file permission: rw-rw-rw-
	return ioutil.WriteFile(secretsPath, []byte(secrets), filePerm)
}

// RemoveSecretsFile removes the secrets file from the user's machine
func RemoveSecretsFile(secretsPath string) error {
	err := os.Remove(secretsPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error removing secrets file: %s. Manual removal recommended", secretsPath))
	}
	return nil
}
