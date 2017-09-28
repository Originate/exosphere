package application

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	prompt "github.com/kofalt/go-prompt"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	err := validateConfigs(deployConfig)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Setting up AWS account...")
	err = aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log(fmt.Sprintf("Building %s %s...\n\n", deployConfig.AppConfig.Name, deployConfig.AppConfig.Version))
	initializer, err := NewInitializer(
		deployConfig.AppConfig,
		deployConfig.Logger,
		deployConfig.AppDir,
		deployConfig.HomeDir,
		deployConfig.DockerComposeProjectName,
		composebuilder.BuildModeDeployProduction)
	if err != nil {
		return err
	}
	err = initializer.Initialize()
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Pushing Docker images to ECR...")
	dockerComposePath := filepath.Join(deployConfig.AppDir, "tmp", "docker-compose.yml")
	imagesMap, err := aws.PushImages(deployConfig, dockerComposePath)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Generating Terraform files...")
	err = terraform.GenerateFile(deployConfig, imagesMap)
	if err != nil {
		return err
	}

	return deployApplication(deployConfig)
}

func validateConfigs(deployConfig types.DeployConfig) error {
	deployConfig.Logger.Log("Validating application configuration...")
	err := deployConfig.AppConfig.Production.ValidateFields()
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Validating service configurations...")
	protectionLevels := deployConfig.AppConfig.GetServiceProtectionLevels()
	serviceData := deployConfig.AppConfig.GetServiceData()
	for serviceName, serviceConfig := range deployConfig.ServiceConfigs {
		err = serviceConfig.Production.ValidateFields(serviceData[serviceName].Location, protectionLevels[serviceName])
		if err != nil {
			return err
		}
	}
	deployConfig.Logger.Log("Validating application dependencies...")
	validatedDependencies := map[string]string{}
	for _, dependency := range deployConfig.AppConfig.Production.Dependencies {
		err = dependency.ValidateFields()
		if err != nil {
			return err
		}
		validatedDependencies[dependency.Name] = dependency.Version
	}

	deployConfig.Logger.Log("Validating service dependencies...")
	for _, serviceConfig := range deployConfig.ServiceConfigs {
		for _, dependency := range serviceConfig.Production.Dependencies {
			if validatedDependencies[dependency.Name] == "" {
				err = dependency.ValidateFields()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func deployApplication(deployConfig types.DeployConfig) error {
	deployConfig.Logger.Log("Retrieving remote state...")
	err := terraform.RunInit(deployConfig)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Retrieving secrets...")
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Planning deployment...")
	err = terraform.RunPlan(deployConfig, secrets)
	if err != nil {
		return err
	}

	if ok := prompt.Confirm("Do you want to apply this plan? (y/n)"); ok {
		deployConfig.Logger.Log("Applying changes...")
		return terraform.RunApply(deployConfig, secrets)
	}
	deployConfig.Logger.Log("Abandoning deployment.")
	return nil
}
