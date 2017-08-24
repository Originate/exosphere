package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// RunInit runs the 'terraform init' command and force copies the remote state
func RunInit(deployConfig types.DeployConfig) error {
	err := util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.LogChannel, "terraform", "init", "-force-copy")
	if err != nil {
		return err
	}
	return err
}

// RunPlan runs the 'terraform plan' command and passes variables in as flags
func RunPlan(deployConfig types.DeployConfig, secrets types.Secrets) error {
	vars, err := CompileVarFlags(deployConfig, secrets)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "plan"}, vars...)
	err = util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.LogChannel, command...)
	if err != nil {
		return err
	}
	return err
}

// RunApply runs the 'terraform apply' command and passes variables in as command flags
func RunApply(deployConfig types.DeployConfig, secrets types.Secrets) error {
	vars, err := CompileVarFlags(deployConfig, secrets)
	if err != nil {
		return err
	}
	command := append([]string{"terraform", "apply"}, vars...)
	err = util.RunAndLog(deployConfig.TerraformDir, []string{}, deployConfig.LogChannel, command...)
	if err != nil {
		return err
	}
	return err
}

// CompileVarFlags compiles the variable flags passed into a Terraform command
func CompileVarFlags(deployConfig types.DeployConfig, secrets types.Secrets) ([]string, error) {
	vars := compileSecrets(secrets)
	envVars, err := compileEnvVars(deployConfig, secrets)
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot compile environment variables")
	}
	vars = append(vars, envVars...)
	return vars, nil
}

func compileSecrets(secrets types.Secrets) []string {
	vars := []string{}
	for k, v := range secrets {
		vars = append(vars, "-var", fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}

func compileEnvVars(deployConfig types.DeployConfig, secrets types.Secrets) ([]string, error) {
	envVars := []string{}
	for serviceName, serviceConfig := range deployConfig.ServiceConfigs {
		serviceEnvVars, serviceSecrets := serviceConfig.GetEnvVars("production")
		for _, secretKey := range serviceSecrets {
			serviceEnvVars[secretKey] = secrets[secretKey]
		}
		serviceEnvVarsStr, err := createEnvVarString(serviceEnvVars)
		if err != nil {
			return []string{}, err
		}
		envVars = append(envVars, "-var", fmt.Sprintf("%s_env_vars=%s", serviceName, serviceEnvVarsStr))
	}
	return envVars, nil
}

func createEnvVarString(envVars map[string]string) (string, error) {
	terraformEnvVars := []map[string]string{}
	for k, v := range envVars {
		envVarPair := map[string]string{
			"key":   k,
			"value": v,
		}
		terraformEnvVars = append(terraformEnvVars, envVarPair)
	}
	envVarsJSON, err := json.Marshal(terraformEnvVars)
	if err != nil {
		return "", err
	}
	envVarsEscaped, err := json.Marshal(string(envVarsJSON))
	if err != nil {
		return "", err
	}
	return string(envVarsEscaped), nil
}
