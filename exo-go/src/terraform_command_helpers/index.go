package terraformCommandHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
)

// TerraformInit runs the 'terraform init' command
func TerraformInit(terraformDir string, log func(string)) error {
	err := processHelpers.RunAndLog(terraformDir, log, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// TerraformPlan runs the 'terraform plan' command
func TerraformPlan(terraformDir string, log func(string)) error {
	err := processHelpers.RunAndLog(terraformDir, log, "terraform", "plan", "-var-file=secret.tfvars")
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
