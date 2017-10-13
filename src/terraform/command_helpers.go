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
	envVars, err := compileServiceEnvVars(deployConfig, secrets)
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot compile environment variables")
	}
	vars = append(vars, envVars...)
	dependencyVars, err := compileDependencyVars(deployConfig)
	if err != nil {
		return []string{}, errors.Wrap(err, "cannot compile dependency variables")
	}
	vars = append(vars, dependencyVars...)
	return append(vars, "-var", fmt.Sprintf("aws_profile=%s", deployConfig.AwsConfig.Profile)), nil
}

// compile all secrets into var flags
func compileSecrets(secrets types.Secrets) []string {
	vars := []string{}
	for k, v := range secrets {
		vars = append(vars, "-var", fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}

// compile docker image var flags for each service
func compileDockerImageVars(deployConfig types.DeployConfig, imagesMap map[string]string) []string {
	vars := []string{}
	for serviceRole := range deployConfig.ServiceConfigs {
		vars = append(vars, "-var", fmt.Sprintf("%s_docker_image=%s", serviceRole, imagesMap[serviceRole]))
	}
	return vars
}

// compile var flags needed for each dependency
func compileDependencyVars(deployConfig types.DeployConfig) ([]string, error) {
	vars := []string{}
	for dependencyName, dependency := range config.GetBuiltAppProductionDependencies(deployConfig.AppConfig, deployConfig.AppDir) {
		varMap, err := dependency.GetDeploymentVariables()
		if err != nil {
			return []string{}, err
		}
		stringifiedVar, err := createEnvVarString(varMap)
		if err != nil {
			return []string{}, err
		}
		vars = append(vars, "-var", fmt.Sprintf("%s_env_vars=%s", dependencyName, stringifiedVar))
	}
	for _, serviceConfig := range deployConfig.ServiceConfigs {
		for dependencyName, dependency := range config.GetBuiltServiceProductionDependencies(serviceConfig, deployConfig.AppConfig, deployConfig.AppDir) {
			varMap, err := dependency.GetDeploymentVariables()
			if err != nil {
				return []string{}, err
			}
			stringifiedVar, err := createEnvVarString(varMap)
			if err != nil {
				return []string{}, err
			}
			vars = append(vars, "-var", fmt.Sprintf("%s_env_vars=%s", dependencyName, stringifiedVar))
		}
	}
	return vars, nil
}

// compile env vars needed for each service
func compileServiceEnvVars(deployConfig types.DeployConfig, secrets types.Secrets) ([]string, error) {
	envVars := []string{}
	for serviceRole, serviceConfig := range deployConfig.ServiceConfigs {
		serviceEnvVars := map[string]string{"ROLE": serviceRole}
		dependencyEnvVars := getDependencyServiceEnvVars(deployConfig, serviceConfig, secrets)
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
		envVars = append(envVars, "-var", fmt.Sprintf("%s_env_vars=%s", serviceRole, serviceEnvVarsStr))
	}
	return envVars, nil
}

// convert an env var key pair in the format of a task definition
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

// get all env vars that a service requires for the its listed dependency
func getDependencyServiceEnvVars(deployConfig types.DeployConfig, serviceConfig types.ServiceConfig, secrets types.Secrets) map[string]string {
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
