package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig deploy.Config) error {
	backendConfig := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	return util.RunAndPipe(deployConfig.GetTerraformDir(), []string{}, deployConfig.Writer, "terraform", "init", "-force-copy", backendConfig)
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string, autoApprove bool) error {
	err := GenerateVarFile(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := []string{"terraform", "apply", fmt.Sprintf("-var-file=%s", terraformVarFile)}
	if autoApprove {
		command = append(command, "-auto-approve")
	}
	return util.RunAndPipe(deployConfig.GetTerraformDir(), []string{}, deployConfig.Writer, command...)
}
