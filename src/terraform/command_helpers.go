package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// CompileVarFlags compiles the variable flags passed into a Terraform command
func CompileVarFlags(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string) ([]string, error) {
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
func compileDockerImageVars(deployConfig deploy.Config, imagesMap map[string]string) []string {
	vars := []string{}
	for serviceRole, serviceConfig := range deployConfig.ServiceConfigs {
		vars = append(vars, "-var", fmt.Sprintf("%s_docker_image=%s", serviceRole, imagesMap[serviceRole]))
		for _, dependency := range serviceConfig.Remote.Dependencies {
			vars = append(vars, "-var", fmt.Sprintf("%s_docker_image=%s", dependency.Name, imagesMap[dependency.Name]))
		}
	}
	for _, dependency := range deployConfig.AppContext.Config.Remote.Dependencies {
		vars = append(vars, "-var", fmt.Sprintf("%s_docker_image=%s", dependency.Name, imagesMap[dependency.Name]))
	}
	return vars
}

// compile var flags needed for each dependency
func compileDependencyVars(deployConfig deploy.Config) ([]string, error) {
	vars := []string{}
	for dependencyName, dependency := range config.GetBuiltRemoteAppDependencies(deployConfig.AppContext) {
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
		for dependencyName, dependency := range config.GetBuiltRemoteServiceDependencies(serviceConfig, deployConfig.AppContext) {
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
func compileServiceEnvVars(deployConfig deploy.Config, secrets types.Secrets) ([]string, error) {
	envVars := []string{}
	for serviceRole, serviceConfig := range deployConfig.ServiceConfigs {
		serviceEnvVars := map[string]string{"ROLE": serviceRole}
		dependencyEnvVars := getDependencyServiceEnvVars(deployConfig, serviceConfig, secrets)
		util.Merge(serviceEnvVars, dependencyEnvVars)
		productionEnvVar, serviceSecrets := serviceConfig.GetEnvVars("remote")
		util.Merge(serviceEnvVars, productionEnvVar)
		for _, secretKey := range serviceSecrets {
			serviceEnvVars[secretKey] = secrets[secretKey]
		}
		endpointEnvVars := getEndpointEnvVars(deployConfig, serviceRole, serviceConfig)
		util.Merge(serviceEnvVars, endpointEnvVars)
		serviceEnvVarsStr, err := createEnvVarString(serviceEnvVars)
		if err != nil {
			return []string{}, err
		}
		envVars = append(envVars, "-var", fmt.Sprintf("%s_env_vars=%s", serviceRole, serviceEnvVarsStr))
	}
	return envVars, nil
}

func getEndpointEnvVars(deployConfig deploy.Config, serviceRole string, serviceConfig types.ServiceConfig) map[string]string {
	s := types.NewServiceEndpoint(deployConfig.AppContext, serviceRole, serviceConfig, nil, deployConfig.BuildMode)
	return s.GetEndpointMappings()
}

// convert an env var key pair in the format of a task definition
// marshals a map[string]string object twice so that it can be escaped properly
// and passed as a command line flag, then properly decoded in Terraform
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
func getDependencyServiceEnvVars(deployConfig deploy.Config, serviceConfig types.ServiceConfig, secrets types.Secrets) map[string]string {
	result := map[string]string{}
	for _, dependency := range config.GetBuiltRemoteAppDependencies(deployConfig.AppContext) {
		util.Merge(
			result,
			dependency.GetDeploymentServiceEnvVariables(secrets),
		)
	}
	for _, dependency := range config.GetBuiltRemoteServiceDependencies(serviceConfig, deployConfig.AppContext) {
		util.Merge(
			result,
			dependency.GetDeploymentServiceEnvVariables(secrets),
		)
	}
	return result
}
