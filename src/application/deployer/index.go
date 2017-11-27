package deployer

import (
	"fmt"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	err := validateConfigs(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Setting up AWS account...")
	err = aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	err = application.GenerateComposeFiles(deployConfig.AppContext)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Pushing Docker images to ECR...")
	imagesMap, err := PushApplicationImages(deployConfig)
	if err != nil {
		return err
	}

	prevTerraformFileContents, err := terraform.ReadTerraformFile(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Generating Terraform files...")
	err = terraform.GenerateFile(deployConfig)
	if err != nil {
		return err
	}

	if deployConfig.DeployServicesOnly {
		err = terraform.CheckTerraformFile(deployConfig, prevTerraformFileContents)
		if err != nil {
			return err
		}
	}

	return deployApplication(deployConfig, imagesMap)
}

func validateConfigs(deployConfig types.DeployConfig) error {
	fmt.Fprintln(deployConfig.Writer, "Validating application configuration...")
	err := deployConfig.AppContext.Config.Production.ValidateFields()
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Validating service configurations...")
	serviceData := deployConfig.AppContext.Config.Services
	for serviceRole, serviceConfig := range deployConfig.ServiceConfigs {
		err = serviceConfig.Production.ValidateFields(serviceData[serviceRole].Location, serviceConfig.Type)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(deployConfig.Writer, "Validating application dependencies...")
	validatedDependencies := map[string]string{}
	for _, dependency := range deployConfig.AppContext.Config.Production.Dependencies {
		err = dependency.ValidateFields()
		if err != nil {
			return err
		}
		validatedDependencies[dependency.Name] = dependency.Version
	}

	fmt.Fprintln(deployConfig.Writer, "Validating service dependencies...")
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

func deployApplication(deployConfig types.DeployConfig, imagesMap map[string]string) error {
	fmt.Fprintln(deployConfig.Writer, "Retrieving remote state...")
	err := terraform.RunInit(deployConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Retrieving secrets...")
	secrets, err := aws.ReadSecrets(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	fmt.Fprintln(deployConfig.Writer, "Applying changes...")
	return terraform.RunApply(deployConfig, secrets, imagesMap, deployConfig.DeployServicesOnly)
}
