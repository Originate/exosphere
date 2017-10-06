package application

import (
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

	deployConfig.Logger.Logf("Building %s %s...", deployConfig.AppConfig.Name, deployConfig.AppConfig.Version)
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

	prevTerraformFileContents, err := terraform.ReadTerraformFile(deployConfig)
	if err != nil {
		return err
	}

	deployConfig.Logger.Log("Generating Terraform files...")
	err = terraform.GenerateFile(deployConfig, imagesMap)
	if err != nil {
		return err
	}

	if deployConfig.DeployServicesOnly {
		err = terraform.CheckTerraformFile(deployConfig, prevTerraformFileContents)
		if err != nil {
			return err
		}
	}

	return deployApplication(deployConfig, imagesMap, prevTerraformFileContents)
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

func deployApplication(deployConfig types.DeployConfig, imagesMap map[string]string, prevTerraformFileContents []byte) error {
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

	applyPlan := true
	if !deployConfig.DeployServicesOnly {
		deployConfig.Logger.Log("Planning deployment...")
		err = terraform.RunPlan(deployConfig, secrets, imagesMap)
		if err != nil {
			return err
		}
		applyPlan = prompt.Confirm("Do you want to apply this plan? (y/n)")
	}

	if applyPlan {
		deployConfig.Logger.Log("Applying changes...")
		return terraform.RunApply(deployConfig, secrets, imagesMap)
	}
	deployConfig.Logger.Log("Abandoning deployment...reverting 'terraform/main.tf' file.")
	return terraform.WriteTerraformFile(string(prevTerraformFileContents), deployConfig.TerraformDir)
}
