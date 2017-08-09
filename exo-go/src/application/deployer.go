package application

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// StartDeploy starts the deployment process
func StartDeploy(deployConfig types.DeployConfig) error {
	terraformDir := filepath.Join(deployConfig.AppDir, "terraform")
	fmt.Printf("Setting up AWS account...\n\n")
	err := aws.InitAccount(deployConfig.AwsConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nGenerating Terraform files...\n\n")
	err = terraform.GenerateFile(deployConfig, terraformDir)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving remote state...\n\n")
	err = terraform.RunInit(terraformDir, deployConfig.Logger)
	if err != nil {
		return err
	}

	secretsPath := filepath.Join(terraformDir, "secrets.tfvars")
	fmt.Printf("\n\nPlanning deployment...\n\n")
	return terraform.RunPlan(terraformDir, secretsPath, deployConfig.Logger)
}
