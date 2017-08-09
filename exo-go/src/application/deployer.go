package application

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// Deployer contains information needed to deploy the application
type Deployer struct {
	AppConfig      types.AppConfig
	ServiceConfigs map[string]types.ServiceConfig
	Logger         chan string
	AppDir         string
	HomeDir        string
}

// Start starts the deployment process
func (d *Deployer) Start() error {
	terraformDir := d.getTerraformDir()
	terraformConfig := types.TerraformConfig{
		AppConfig:      d.AppConfig,
		ServiceConfigs: d.ServiceConfigs,
		AppDir:         d.AppDir,
		HomeDir:        d.HomeDir,
		TerraformDir:   terraformDir,
		RemoteBucket:   fmt.Sprintf("%s-terraform", d.AppConfig.Name),
		LockTable:      "TerraformLocks",
		Region:         "us-west-2", //TODO prompt user for this
	}

	fmt.Printf("Setting up AWS account...\n\n")
	err := aws.InitAccount(terraformConfig.RemoteBucket, terraformConfig.LockTable, terraformConfig.Region)
	if err != nil {
		return err
	}

	fmt.Printf("Building %s %s...\n\n", d.AppConfig.Name, d.AppConfig.Version)
	initializer, err := NewInitializer(d.AppConfig, d.Logger, "exo-deploy", d.AppDir, d.HomeDir)
	if err != nil {
		return err
	}
	err = initializer.Initialize()
	if err != nil {
		return err
	}

	fmt.Printf("\n\nPushing Docker images to ECR...\n\n")
	dockerComposePath := filepath.Join(d.AppDir, "tmp", "docker-compose.yml")
	err = aws.PushImages(d.AppConfig, dockerComposePath, terraformConfig.Region)
	if err != nil {
		return err
	}

	err = terraform.GenerateFile(terraformConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRetrieving remote state...\n\n")
	err = terraform.RunInit(terraformDir, d.Logger)
	if err != nil {
		return err
	}

	secretsFile := d.getSecretsFile(d.AppDir)
	fmt.Printf("\n\nPlanning deployment...\n\n")
	return terraform.RunPlan(terraformDir, secretsFile, d.Logger)
}

func (d *Deployer) getTerraformDir() string {
	return filepath.Join(d.AppDir, "terraform")
}

func (d *Deployer) getSecretsFile(appDir string) string {
	return filepath.Join(d.getTerraformDir(), "secret.tfvars")
}
