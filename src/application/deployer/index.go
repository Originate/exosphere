package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
)

func setupDeploy(deployConfig deploy.Config) (types.Secrets, error) {
	err := ValidateConfigs(deployConfig)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Validating Terraform files...")
	if err != nil {
		return nil, err
	}
	err = terraform.GenerateCheck(deployConfig)
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Setting up AWS account...")
	if err != nil {
		return nil, err
	}
	err = aws.InitAccount(deployConfig.GetAwsOptions())
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(deployConfig.Writer, "Retrieving secrets...")
	return aws.ReadSecrets(deployConfig.GetAwsOptions())
}

// DeployInfrastructure deploys the infrastructure for the application
func DeployInfrastructure(deployConfig deploy.Config) error {
	secrets, err := setupDeploy(deployConfig)
	if err != nil {
		return err
	}
	terraformDir := deployConfig.GetInfrastructureTerraformDir()
	err = terraform.RunInit(deployConfig, terraformDir)
	if err != nil {
		return err
	}
	err = terraform.GenerateInfrastructureVarFile(deployConfig, secrets)
	if err != nil {
		return err
	}
	return terraform.RunApply(deployConfig, terraformDir, deployConfig.AutoApprove)
}

// DeployServices deploys the services for the application
func DeployServices(deployConfig deploy.Config) error {
	secrets, err := setupDeploy(deployConfig)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Pushing Docker images to ECR...")
	if err != nil {
		return err
	}
	serviceDockerImagesMap, err := PushApplicationImages(deployConfig)
	if err != nil {
		return err
	}
	terraformDir := deployConfig.GetServicesTerraformDir()
	err = terraform.RunInit(deployConfig, terraformDir)
	if err != nil {
		return err
	}
	err = terraform.GenerateServicesVarFile(deployConfig, secrets, serviceDockerImagesMap)
	if err != nil {
		return err
	}
	return terraform.RunApply(deployConfig, terraformDir, deployConfig.AutoApprove)
}

// ValidateConfigs validates application/service deployment configuration fields
func ValidateConfigs(deployConfig deploy.Config) error {
	fmt.Fprintln(deployConfig.Writer, "Validating application configuration...")
	err := deployConfig.AppContext.Config.Remote.ValidateFields(deployConfig.RemoteEnvironmentID)
	if err != nil {
		return err
	}
	fmt.Fprintln(deployConfig.Writer, "Validating service configurations...")
	for _, serviceContext := range deployConfig.AppContext.ServiceContexts {
		err = serviceContext.Config.ValidateDeployFields(serviceContext.Source.Location, deployConfig.RemoteEnvironmentID)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(deployConfig.Writer, "Validating application dependencies...")
	for _, dependency := range deployConfig.AppContext.Config.Remote.Dependencies {
		err = dependency.ValidateFields()
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(deployConfig.Writer, "Validating service dependencies...")
	for _, serviceContext := range deployConfig.AppContext.ServiceContexts {
		for _, dependency := range serviceContext.Config.Remote.Dependencies {
			err = dependency.ValidateFields()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
