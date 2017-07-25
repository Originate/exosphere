package terraformCommandHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
)

// TerraformInit runs the 'terraform init' command
func TerraformInit(terraformDir string, log func(string)) error {
	return processHelpers.RunAndLog(terraformDir, log, "terraform", "init", "-force-copy")
}

// TerraformPlan runs the 'terraform plan' command
func TerraformPlan(terraformDir string, log func(string)) error {
	return processHelpers.RunAndLog(terraformDir, log, "terraform", "plan", "-var-file=secret.tfvars")
}
