package application

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	deployConfig.LogChannel <- "Setting up AWS account..."
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- fmt.Sprintf("Building %s %s...\n\n", deployConfig.AppConfig.Name, deployConfig.AppConfig.Version)
	initializer, err := NewInitializer(deployConfig.AppConfig, deployConfig.LogChannel, "exo-deploy", deployConfig.AppDir, deployConfig.HomeDir)
	if err != nil {
		return err
	}
	err = initializer.Initialize()
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Pushing Docker images to ECR..."
	dockerComposePath := filepath.Join(deployConfig.AppDir, "tmp", "docker-compose.yml")
	imagesMap, err := aws.PushImages(deployConfig, dockerComposePath)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Generating Terraform files..."
	err = terraform.GenerateFile(deployConfig, imagesMap)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Retrieving remote state..."
	err = terraform.RunInit(deployConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Retrieving secrets..."
	err = writeSecretsFile(deployConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Planning deployment..."
	return terraform.RunPlan(deployConfig)
}

func writeSecretsFile(deployConfig types.DeployConfig) error {
	return nil //TODO in next PR
}

// RemoveSecretsFile removes the secrets file from the user's machine
func RemoveSecretsFile(secretsPath string) error {
	if util.DoesFileExist(secretsPath) {
		err := os.Remove(secretsPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error removing secrets file: %s. Manual removal recommended", secretsPath))
		}
	}
	return nil
}
