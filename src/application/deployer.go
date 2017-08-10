package application

import (
	"path/filepath"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	terraformDir := filepath.Join(deployConfig.AppDir, "terraform")
	deployConfig.LogChannel <- "Setting up AWS account..."
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Generating Terraform files..."
	err = terraform.GenerateFile(deployConfig, terraformDir)
	if err != nil {
		return err
	}

	deployConfig.LogChannel <- "Retrieving remote state..."
	err = terraform.RunInit(terraformDir, deployConfig.LogChannel)
	if err != nil {
		return err
	}

	secretsPath := filepath.Join(terraformDir, "secrets.tfvars")
	deployConfig.LogChannel <- "Planning deployment..."
	return terraform.RunPlan(terraformDir, secretsPath, deployConfig.LogChannel)
}
