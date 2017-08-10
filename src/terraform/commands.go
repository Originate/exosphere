package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig types.DeployConfig) error {
	err := util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.Logger, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// RunPlan runs the 'terraform plan' command and points to a secrets file
func RunPlan(deployConfig types.DeployConfig) error {
	err := util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.Logger, "terraform", "plan", fmt.Sprintf("-var-file=%s", deployConfig.SecretsPath))
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
