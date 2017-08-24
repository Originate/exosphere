package terraform

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig types.DeployConfig) error {
	err := util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.LogChannel, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// RunPlan runs the 'terraform plan' command and passes variables in as flags
func RunPlan(deployConfig types.DeployConfig, secrets types.Secrets) error {
	command := []string{"terraform", "plan"}
	vars, err := CompileVarFlags(deployConfig, secrets)
	if err != nil {
		return err
	}
	command = append(command, vars...)
	err = util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.LogChannel, command...)
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
