package terraformCommandHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
)

func TerraformInit(terraformDir string, logger *logger.Logger) error {
	return processHelpers.RunAndLog(terraformDir, logger, "terraform", "init")
}

func TerraformPlan(terraformDir string, logger *logger.Logger) error {
	return processHelpers.RunAndLog(terraformDir, logger, "terraform", "plan", "-vars-file=secret.tfvars")
}
