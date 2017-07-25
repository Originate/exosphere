package terraformCommandHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
)

// TerraformInit runs the 'terraform init' command
func TerraformInit(terraformDir string, log func(string)) error {
	return processHelpers.RunAndLog(terraformDir, log, "terraform", "init")
}

// TerraformPlan runs the 'terraform plan' command
func TerraformPlan(terraformDir string, log func(string)) error {
	return processHelpers.RunAndLog(terraformDir, log, "terraform", "plan", "-vars-file=secret.tfvars")
}
