package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

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
	dependencyEnvVars := map[string]string{}
	for _, dependency := range deployConfig.AppConfig.Dependencies {
		util.Merge(
			dependencyEnvVars,
			config.NewAppDependency(dependency, deployConfig.AppConfig, deployConfig.AppDir, deployConfig.HomeDir).GetDeploymentServiceEnvVariables(),
		)
	}
	for serviceName, serviceConfig := range deployConfig.ServiceConfigs {
		serviceEnvVars, serviceSecrets := serviceConfig.GetEnvVars("production")
		for _, secretKey := range serviceSecrets {
			serviceEnvVars[secretKey] = secrets[secretKey]
		}
		util.Merge(serviceEnvVars, dependencyEnvVars)
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
			"name":  k,
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
