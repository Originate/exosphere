package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig deploy.Config) error {
	backendConfigProfile := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	backendConfigRegion := fmt.Sprintf("-backend-config=region=%s", deployConfig.AwsConfig.Region)
	return util.RunAndPipe(deployConfig.GetTerraformDir(), []string{}, deployConfig.Writer, "terraform", "init", "-force-copy", backendConfigProfile, backendConfigRegion)
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string, autoApprove bool) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "apply"}, vars...)
	if autoApprove {
		command = append(command, "-auto-approve")
	}
	return util.RunAndPipe(deployConfig.GetTerraformDir(), []string{}, deployConfig.Writer, command...)
}
