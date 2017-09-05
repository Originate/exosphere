package application

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	prompt "github.com/segmentio/go-prompt"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	deployConfig.LogChannel <- "Setting up AWS account..."
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- fmt.Sprintf("Building %s %s...\n\n", deployConfig.AppConfig.Name, deployConfig.AppConfig.Version)
	initializer, err := NewInitializer(
		deployConfig.AppConfig,
		deployConfig.LogChannel,
		"exo-deploy",
		deployConfig.AppDir,
		deployConfig.HomeDir,
		deployConfig.DockerComposeProjectName)
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
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Planning deployment..."
	terraform.RunPlan(deployConfig, secrets)
	if err != nil {
		return err
	}

	if ok := prompt.Confirm("Do you want to apply this plan?"); ok {
		deployConfig.LogChannel <- "Applying changes..."
		return terraform.RunApply(deployConfig, secrets)
	}
	deployConfig.LogChannel <- "Abandoning deployment."
	return nil
}
