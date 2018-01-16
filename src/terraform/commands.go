package terraform

import (
	"fmt"

	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig deploy.Config, terraformDir string) error {
	command := []string{"terraform", "init", "-force-copy"}
	for key, value := range getBackendConfigMap(deployConfig) {
		command = append(command, fmt.Sprintf("-backend-config=%s=%s", key, value))
	}
	return util.RunAndPipe(terraformDir, []string{}, deployConfig.Writer, command...)
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig deploy.Config, terraformDir string, autoApprove bool) error {
	command := []string{"terraform", "apply", fmt.Sprintf("-var-file=%s.tfvars", deployConfig.RemoteEnvironmentID)}
	if autoApprove {
		command = append(command, "-auto-approve")
	}
	return util.RunAndPipe(terraformDir, []string{}, deployConfig.Writer, command...)
}

func getBackendConfigMap(deployConfig deploy.Config) map[string]string {
	return map[string]string{
		"bucket":  deployConfig.GetAwsBucketName(),
		"profile": deployConfig.AwsProfile,
		"region":  deployConfig.GetRemoteEnvironment().Region,
	}
}
