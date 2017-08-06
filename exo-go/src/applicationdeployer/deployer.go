package applicationdeployer

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws"
	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/terraform"
)

// Deployer contains information needed to deploy the application
type Deployer struct {
	AppConfig      config.AppConfig
	ServiceConfigs map[string]config.ServiceConfig
	Logger         chan string
	AppDir         string
	HomeDir        string
}

// Start starts the deployment process
func (d *Deployer) Start() error {
	terraformDir := d.getTerraformDir()
	terraformConfig := config.TerraformConfig{
		AppConfig:      d.AppConfig,
		ServiceConfigs: d.ServiceConfigs,
		AppDir:         d.AppDir,
		HomeDir:        d.HomeDir,
		TerraformDir:   terraformDir,
		RemoteBucket:   fmt.Sprintf("%s-terraform", d.AppConfig.Name),
		LockTable:      "TerraformLocks",
		Region:         "us-west-2", //TODO prompt user for this
	}

	err := aws.InitAccount(terraformConfig.RemoteBucket, terraformConfig.LockTable, terraformConfig.Region)
	if err != nil {
		return err
	}

	err = terraform.GenerateFile(terraformConfig)
	if err != nil {
		return err
	}

	err = terraform.RunInit(terraformDir, d.Logger)
	if err != nil {
		return err
	}

	secretsFile := d.getSecretsFile(d.AppDir)
	return terraform.RunPlan(terraformDir, secretsFile, d.Logger)
}

func (d *Deployer) getTerraformDir() string {
	return filepath.Join(d.AppDir, "terraform")
}

func (d *Deployer) getSecretsFile(appDir string) string {
	return filepath.Join(d.getTerraformDir(), "secret.tfvars")
}
