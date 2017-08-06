package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/runplus"
	"github.com/pkg/errors"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(terraformDir string, logChannel chan string) error {
	err := runplus.AndLog(terraformDir, []string{}, logChannel, "terraform", "init", "-force-copy")
	if err != nil {
		return errors.Wrap(err, "'terraform init' failed")
	}
	return err
}

// RunPlan runs the 'terraform plan' command and points to a secrets file
func RunPlan(terraformDir, secretsFile string, logChannel chan string) error {
	err := runplus.AndLog(terraformDir, []string{}, logChannel, "terraform", "plan", fmt.Sprintf("-var-file=%s", secretsFile))
	if err != nil {
		return errors.Wrap(err, "'terraform plan' failed")
	}
	return err
}
