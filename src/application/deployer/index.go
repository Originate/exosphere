package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
)

// SetupDeploy starts the deployment process
// nolint gocyclo
func SetupDeploy(deployConfig deploy.Config) (map[string]string, types.Secrets, error) {
	err := ValidateConfigs(deployConfig)
	if err != nil {
		return nil, nil, err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Validating Terraform files...")
	err = terraform.GenerateCheck(deployConfig)
	if err != nil {
		return nil, nil, err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Setting up AWS account...")
	if err != nil {
		return nil, nil, err
	}
	err = aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return nil, nil, err
	}
	_, err = fmt.Fprintln(deployConfig.Writer, "Pushing Docker images to ECR...")
	if err != nil {
		return nil, nil, err
	}
	imagesMap, err := PushApplicationImages(deployConfig)
	if err != nil {
		return nil, nil, err
	}
	fmt.Fprintln(deployConfig.Writer, "Retrieving secrets...")
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	return imagesMap, secrets, err
}

func DeployInfrastructure(deployConfig deploy.Config, imagesMap map[string]string, secrets types.Secrets) error {
	fmt.Fprintln(deployConfig.Writer, "Retrieving remote state...")
	terraformDir := deployConfig.GetInfrastructureTerraformDir()
	err := terraform.RunInit(deployConfig, terraformDir)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Applying changes...")
	return terraform.RunApply(deployConfig, secrets, imagesMap, terraformDir, deployConfig.AutoApprove)
}

func DeployServices(deployConfig deploy.Config, imagesMap map[string]string, secrets types.Secrets) error {
	fmt.Fprintln(deployConfig.Writer, "Retrieving remote state...")
	terraformDir := deployConfig.GetServicesTerraformDir()
	err := terraform.RunInit(deployConfig, terraformDir)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Applying changes...")
	return terraform.RunApply(deployConfig, secrets, imagesMap, terraformDir, deployConfig.AutoApprove)
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
