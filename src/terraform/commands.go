package terraform

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig deploy.Config) error {
	homeDir, err := util.GetHomeDirectory()
	if err != nil {
		return err
	}
	backendConfig := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	return tools.RunInDockerContainer(tools.RunConfig{
		Volumes:    []string{fmt.Sprintf("%s:/app", deployConfig.TerraformDir), fmt.Sprintf("%s/.aws:/root/.aws", homeDir)},
		WorkingDir: "/app",
		ImageName:  fmt.Sprintf("%s:%s", TerraformImage, TerraformVersion),
		Command:    []string{"init", "-force-copy", backendConfig},
		Writer:     deployConfig.Writer,
	})
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string, autoApprove bool) error {
	vars, err := CompileVarFlags(deployConfig, secrets, imagesMap)
	if err != nil {
		return err
	}
	command := append([]string{"apply"}, vars...)
	if autoApprove {
		command = append(command, "-auto-approve")
	}
	homeDir, err := util.GetHomeDirectory()
	if err != nil {
		return err
	}
	err = tools.RunInDockerContainer(tools.RunConfig{
		Volumes:     []string{fmt.Sprintf("%s:/app", deployConfig.TerraformDir), fmt.Sprintf("%s/.aws:/root/.aws", homeDir)},
		Interactive: true,
		WorkingDir:  "/app",
		ImageName:   fmt.Sprintf("%s:%s", TerraformImage, TerraformVersion),
		Command:     command,
		Writer:      deployConfig.Writer,
	})
	if err != nil && strings.Contains(err.Error(), "exit status") {
		return nil
	}
	return err
}
