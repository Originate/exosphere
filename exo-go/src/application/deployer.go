package application

import (
	"path/filepath"

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

	err := terraform.GenerateFile(d.AppConfig, d.ServiceConfigs, d.AppDir, d.HomeDir, terraformDir)
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
