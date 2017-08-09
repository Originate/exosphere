package application

import (
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	terraformDir := filepath.Join(deployConfig.AppDir, "terraform")
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	err = terraform.GenerateFile(deployConfig, terraformDir)
	if err != nil {
		return err
	}

	err = terraform.RunInit(terraformDir, deployConfig.Logger)
	if err != nil {
		return err
	}

	secretsPath := filepath.Join(terraformDir, "secrets.tfvars")
	return terraform.RunPlan(terraformDir, secretsPath, deployConfig.Logger)
}
