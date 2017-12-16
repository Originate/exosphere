package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/types/endpoints"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// CompileVarFlags compiles the variable flags passed into a Terraform command
func CompileVarFlags(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string) ([]string, error) {
	varMap, err := GetVarMap(deployConfig, secrets, imagesMap)
	if err != nil {
		return []string{}, err
	}
	return buildFlags(varMap), nil
}

// GetVarMap compiles the variables passed into a Terraform command
func GetVarMap(deployConfig deploy.Config, secrets types.Secrets, imagesMap map[string]string) (map[string]string, error) {
	varMap := getDockerImageVarMap(deployConfig, imagesMap)
	servicesVarMap, err := getServicesVarMap(deployConfig, secrets)
	if err != nil {
		return map[string]string{}, errors.Wrap(err, "cannot compile service environment variables")
	}
	util.Merge(varMap, servicesVarMap)
	dependenciesVarMap, err := getDependenciesVarMap(deployConfig)
	if err != nil {
		return map[string]string{}, errors.Wrap(err, "cannot compile dependency variables")
	}
	util.Merge(varMap, dependenciesVarMap)
	util.Merge(varMap, secrets)
	util.Merge(varMap, getAwsVarMap(deployConfig))
	util.Merge(varMap, getURLVarMap(deployConfig))
	varMap["env"] = deployConfig.RemoteEnvironmentID
	return varMap, nil
}

// buildFlags builds a map[string]string object into terraform var flags
func buildFlags(varMap map[string]string) []string {
	varFlags := []string{}
	for k, v := range varMap {
		varFlags = append(varFlags, "-var", fmt.Sprintf("%s=%s", k, v))
	}
	return varFlags
}

// getDockerImageVarMap compiles the docker image variables for each service
func getDockerImageVarMap(deployConfig deploy.Config, imagesMap map[string]string) map[string]string {
	dockerImages := map[string]string{}
	for serviceRole, serviceContext := range deployConfig.AppContext.ServiceContexts {
		dockerImages[fmt.Sprintf("%s_docker_image", serviceRole)] = imagesMap[serviceRole]
		for dependencyName := range serviceContext.Config.Remote.Dependencies {
			dockerImages[fmt.Sprintf("%s_docker_image", dependencyName)] = imagesMap[dependencyName]
		}
	}
	for dependencyName := range deployConfig.AppContext.Config.Remote.Dependencies {
		dockerImages[fmt.Sprintf("%s_docker_image", dependencyName)] = imagesMap[dependencyName]
	}
	return dockerImages
}

// getDependenciesVarMap compiles variables  needed for each dependency
func getDependenciesVarMap(deployConfig deploy.Config) (map[string]string, error) {
	dependencyVars := map[string]string{}
	for dependencyName, dependency := range config.GetBuiltRemoteAppDependencies(deployConfig.AppContext) {
		varMap, err := dependency.GetDeploymentVariables()
		if err != nil {
			return map[string]string{}, err
		}
		stringifiedVar, err := createEnvVarString(varMap)
		if err != nil {
			return map[string]string{}, err
		}
		dependencyVars[fmt.Sprintf("%s_env_vars", dependencyName)] = stringifiedVar
	}
	for _, serviceContext := range deployConfig.AppContext.ServiceContexts {
		for dependencyName, dependency := range config.GetBuiltRemoteServiceDependencies(serviceContext.Config, deployConfig.AppContext) {
			varMap, err := dependency.GetDeploymentVariables()
			if err != nil {
				return map[string]string{}, err
			}
			stringifiedVar, err := createEnvVarString(varMap)
			if err != nil {
				return map[string]string{}, err
			}
			dependencyVars[fmt.Sprintf("%s_env_vars", dependencyName)] = stringifiedVar
		}
	}
	return dependencyVars, nil
}

// getServicesVarMap compiles env vars needed for each service
func getServicesVarMap(deployConfig deploy.Config, secrets types.Secrets) (map[string]string, error) {
	envVars := map[string]string{}
	serviceEndpoints := endpoints.NewServiceEndpoints(deployConfig.AppContext, types.BuildModeDeploy, deployConfig.RemoteEnvironmentID)
	for serviceRole, serviceContext := range deployConfig.AppContext.ServiceContexts {
		serviceEnvVars := map[string]string{"ROLE": serviceRole}
		util.Merge(serviceEnvVars, deployConfig.AppContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].Environment)
		dependencyEnvVars := getDependencyServiceEnvVars(deployConfig, serviceContext.Config, secrets)
		util.Merge(serviceEnvVars, dependencyEnvVars)
		remoteEnvironment := serviceContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID]
		util.Merge(serviceEnvVars, remoteEnvironment.Environment)
		for _, secretKey := range remoteEnvironment.Secrets {
			serviceEnvVars[secretKey] = secrets[secretKey]
		}
		for _, secretKey := range deployConfig.AppContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].Secrets {
			serviceEnvVars[secretKey] = secrets[secretKey]
		}
		endpointEnvVars := serviceEndpoints.GetServiceEndpointEnvVars(serviceRole)
		util.Merge(serviceEnvVars, endpointEnvVars)
		serviceEnvVarsStr, err := createEnvVarString(serviceEnvVars)
		if err != nil {
			return map[string]string{}, err
		}
		envVars[fmt.Sprintf("%s_env_vars", serviceRole)] = serviceEnvVarsStr
	}
	return envVars, nil
}

func getAwsVarMap(deployConfig deploy.Config) map[string]string {
	return map[string]string{
		"aws_profile":             deployConfig.AwsConfig.Profile,
		"aws_region":              deployConfig.AwsConfig.Region,
		"aws_account_id":          deployConfig.AwsConfig.AccountID,
		"aws_ssl_certificate_arn": deployConfig.AwsConfig.SslCertificateArn,
	}
}

func getURLVarMap(deployConfig deploy.Config) map[string]string {
	varMap := map[string]string{
		"application_url": deployConfig.AppContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].URL,
	}
	for serviceRole, serviceContext := range deployConfig.AppContext.ServiceContexts {
		if serviceContext.Config.Type == types.ServiceTypePublic {
			varMap[fmt.Sprintf("%s_url", serviceRole)] = serviceContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].URL
		}
	}
	return varMap
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
