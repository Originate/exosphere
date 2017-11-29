package terraform

import (
	"fmt"

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
	volumes := []string{
		fmt.Sprintf("%s:/app", deployConfig.TerraformDir),
		fmt.Sprintf("%s/.aws:/root/.aws", homeDir),
	}
	backendConfig := fmt.Sprintf("-backend-config=profile=%s", deployConfig.AwsConfig.Profile)
	dockerRunConfig := tools.RunConfig{
		Volumes:    volumes,
		WorkingDir: "/app",
		ImageName:  fmt.Sprintf("%s:%s", TerraformImage, TerraformVersion),
		Command:    []string{"init", "-force-copy", backendConfig},
		Writer:     deployConfig.Writer,
	}
	return tools.RunInDockerContainer(dockerRunConfig)
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
	volumes := []string{
		fmt.Sprintf("%s:/app", deployConfig.TerraformDir),
		fmt.Sprintf("%s/.aws:/root/.aws", homeDir),
	}
	dockerRunConfig := tools.RunConfig{
		Volumes:     volumes,
		Interactive: true,
		WorkingDir:  "/app",
		ImageName:   fmt.Sprintf("%s:%s", TerraformImage, TerraformVersion),
		Command:     command,
		Writer:      deployConfig.Writer,
	}
	return tools.RunInDockerContainer(dockerRunConfig)
}
