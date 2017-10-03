package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig types.DeployConfig) error {
	backendConfig := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	return util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.Logger, "terraform", "init", "-force-copy", backendConfig)
}

// RunPlan runs the 'terraform plan' command and passes variables in as flags
func RunPlan(deployConfig types.DeployConfig, secrets types.Secrets, imagesMap map[string]string) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "plan"}, vars...)
	return util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.Logger, command...)
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig types.DeployConfig, secrets types.Secrets, imagesMap map[string]string) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "apply"}, vars...)
	return util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.Logger, command...)
}
