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
func CompileVarFlags(deployConfig types.DeployConfig, secrets types.Secrets, imagesMap map[string]string) ([]string, error) {
	vars := compileSecrets(secrets)
	imageVars := compileDockerImageVars(deployConfig, imagesMap)
	vars = append(vars, imageVars...)
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

func compileDockerImageVars(deployConfig types.DeployConfig, imagesMap map[string]string) []string {
	vars := []string{}
	for serviceName := range deployConfig.ServiceConfigs {
		vars = append(vars, "-var", fmt.Sprintf("%s_docker_image=%s", serviceName, imagesMap[serviceName]))
	}
	return vars
}

func compileEnvVars(deployConfig types.DeployConfig, secrets types.Secrets) ([]string, error) {
	envVars := []string{}
	for serviceName, serviceConfig := range deployConfig.ServiceConfigs {
		serviceEnvVars := map[string]string{"ROLE": serviceName}
		dependencyEnvVars := getDependencyEnvVars(deployConfig, serviceConfig, secrets)
		util.Merge(serviceEnvVars, dependencyEnvVars)
		productionEnvVar, serviceSecrets := serviceConfig.GetEnvVars("production")
		util.Merge(serviceEnvVars, productionEnvVar)
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

func getDependencyEnvVars(deployConfig types.DeployConfig, serviceConfig types.ServiceConfig, secrets types.Secrets) map[string]string {
	result := map[string]string{}
	for _, dependency := range config.GetBuiltAppProductionDependencies(deployConfig.AppConfig, deployConfig.AppDir) {
		util.Merge(
			result,
			dependency.GetDeploymentServiceEnvVariables(secrets),
		)
	}
	for _, dependency := range config.GetBuiltServiceProductionDependencies(serviceConfig, deployConfig.AppConfig, deployConfig.AppDir) {
		util.Merge(
			result,
			dependency.GetDeploymentServiceEnvVariables(secrets),
		)
	}
	return result
}
