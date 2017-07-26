package appDeployer

import (
	"fmt"
	"path/filepath"

	"github.com/Originate/exosphere/exo-go/src/aws_helpers"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/terraform_command_helpers"
	"github.com/Originate/exosphere/exo-go/src/terraform_file_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppDeployer contains information needed to deploy the application
type AppDeployer struct {
	AppConfig      types.AppConfig
	ServiceConfigs map[string]types.ServiceConfig
	AppDir         string
	HomeDir        string
	Logger         *logger.Logger
}

// NewAppDeployer is the constructor for the appDeployer class
func NewAppDeployer(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir, homeDir string, logger *logger.Logger) *AppDeployer {
	return &AppDeployer{
		AppConfig:      appConfig,
		ServiceConfigs: serviceConfigs,
		AppDir:         appDir,
		HomeDir:        homeDir,
		Logger:         logger,
	}
}

// Start starts the deployment process
func (d *AppDeployer) Start() error {
	terraformDir := getTerraformDir(d.AppDir)
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

	// Configure blank AWS account
	err := awsHelper.Configure(terraformConfig.RemoteBucket, terraformConfig.LockTable, terraformConfig.Region)
	if err != nil {
		return err
	}

	// Generate Terraform Scripts
	err = terraformFileHelpers.GenerateTerraformFile(terraformConfig)
	if err != nil {
		return err
	}

	// Run Terraform commands
	err = terraformCommandHelpers.TerraformInit(terraformDir, d.write)
	if err != nil {
		return err
	}

	secretsFile := getSecretsFile(d.AppDir)
	return terraformCommandHelpers.TerraformPlan(terraformDir, secretsFile, d.write)
}

func (d *AppDeployer) write(text string) {
	err := d.Logger.Log("exo-deploy", text, true)
	if err != nil {
		fmt.Printf("Error logging exo-deploy output: %v\n", err)
	}
}

func getTerraformDir(appDir string) string {
	return filepath.Join(appDir, "terraform")
}

func getSecretsFile(appDir string) string {
	return filepath.Join(appDir, "terraform", "secret.tfvars")
}
