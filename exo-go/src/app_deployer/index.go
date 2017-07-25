package appDeployer

import (
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppDeployer contains information needed to deploy the application
type AppDeployer struct {
	AppConfig      types.AppConfig
	ServiceConfigs map[string]types.ServiceConfig
	Logger         *logger.Logger
	AppDir         string
	HomeDir        string
}

// NewAppDeployer is the constructor for the appDeployer class
func NewAppDeployer(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir, homeDir string) *AppDeployer {
	return &AppDeployer{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
	}
}

// Start starts the deployment process
func (d *AppDeployer) Start() error {
	terraformDir := getTerraformDir(d.AppDir)
	return terraformFileHelpers.GenerateTerraformFile(d.AppConfig, d.ServiceConfigs, d.AppDir, d.HomeDir, terraformDir)
}

func getTerraformDir(appDir string) string {
	return filepath.Join(appDir, "terraform")
}
