package terraformCommandHelpers

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/process_helpers"
	"github.com/pkg/errors"
)

// TerraformInit runs the 'terraform init' command and force copies the remote state
func TerraformInit(terraformDir string, log func(string)) error {
	err := processHelpers.RunAndLog(terraformDir, log, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// TerraformPlan runs the 'terraform plan' command and points to a secrets file
func TerraformPlan(terraformDir, secretsFile string, log func(string)) error {
	err := processHelpers.RunAndLog(terraformDir, log, "terraform", "plan", fmt.Sprintf("-var-file=%s", secretsFile))
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
