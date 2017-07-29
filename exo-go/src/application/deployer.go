package application

import (
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/terraform"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// Deployer contains information needed to deploy the application
type Deployer struct {
	AppConfig      types.AppConfig
	ServiceConfigs map[string]types.ServiceConfig
	Logger         *logger.Logger
	AppDir         string
	HomeDir        string
}

// NewDeployer is the constructor for the Deployer class
func NewDeployer(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir, homeDir string) *Deployer {
	return &Deployer{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
	}
}

// Start starts the deployment process
func (d *Deployer) Start() error {
	terraformDir := getTerraformDir(d.AppDir)
	return terraform.GenerateFile(d.AppConfig, d.ServiceConfigs, d.AppDir, d.HomeDir, terraformDir)
}

func getTerraformDir(appDir string) string {
	return filepath.Join(appDir, "terraform")
}
