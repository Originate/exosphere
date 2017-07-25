package appDeployer

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/terraform_command_helpers"
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
	TerraformDir   string
}

// NewAppDeployer is the constructor for the appDeployer class
func NewAppDeployer(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir, homeDir string, logger *logger.Logger) *AppDeployer {
	return &AppDeployer{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
		Logger:         logger,
		TerraformDir:   filepath.Join(appDir, "terraform"),
	}
}

// Start starts the deployment process
func (d *AppDeployer) Start() error {
	err := terraformFileHelpers.GenerateTerraformFile(d.AppConfig, d.ServiceConfigs, d.AppDir, d.HomeDir, d.TerraformDir)
	if err != nil {
		return err
	}
	return terraformCommandHelpers.TerraformInit(d.TerraformDir, d.write)
}

func (d *AppDeployer) write(text string) {
	err := d.Logger.Log("exo-deploy", text, true)
	if err != nil {
		fmt.Printf("Error logging exo-deploy output: %v\n", err)
	}
}
