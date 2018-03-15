package terraform

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
)

// GenerateCustomTerraformVarFile compiles the variable flags passed into a Terraform command
func GenerateCustomTerraformVarFile(deployConfig deploy.Config, secrets types.Secrets) error {
	varMap := getSharedVarMap(deployConfig, secrets)
	return writeVarFile(varMap, deployConfig.GetCustomTerraformDir(), deployConfig.RemoteEnvironmentID)
}
