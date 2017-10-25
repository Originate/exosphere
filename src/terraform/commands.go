package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig types.DeployConfig) error {
	backendConfig := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	return util.RunAndPipe(deployConfig.TerraformDir, []string{}, deployConfig.Writer, "terraform", "init", "-force-copy", backendConfig)
}

// RunPlan runs the 'terraform plan' command and passes variables in as flags
func RunPlan(deployConfig types.DeployConfig, secrets types.Secrets, imagesMap map[string]string) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "plan"}, vars...)
	return util.RunAndPipe(deployConfig.TerraformDir, []string{}, deployConfig.Writer, command...)
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig types.DeployConfig, secrets types.Secrets, imagesMap map[string]string) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "apply"}, vars...)
	return util.RunAndPipe(deployConfig.TerraformDir, []string{}, deployConfig.Writer, command...)
}
