package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/pkg/errors"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(terraformDir string, logChannel chan string) error {
	err := util.RunAndLog(terraformDir, []string{}, logChannel, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// RunPlan runs the 'terraform plan' command and points to a secrets file
func RunPlan(terraformDir, secretsPath string, logChannel chan string) error {
	err := util.RunAndLog(terraformDir, []string{}, logChannel, "terraform", "plan", fmt.Sprintf("-var-file=%s", secretsPath))
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
