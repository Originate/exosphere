package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/types/endpoints"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// GenerateServicesVarFile compiles the variable flags passed into a Terraform command
func GenerateServicesVarFile(deployConfig deploy.Config, secrets types.Secrets, serviceDockerImagesMap map[string]string) error {
	varMap, err := GetServicesVarMap(deployConfig, secrets, serviceDockerImagesMap)
	if err != nil {
		return err
	}
	jsonVarMap, err := json.MarshalIndent(varMap, "", "  ")
	if err != nil {
		return err
	}
	return WriteToTerraformDir(string(jsonVarMap), fmt.Sprintf("%s.tfvars", deployConfig.RemoteEnvironmentID), deployConfig.GetServicesTerraformDir())
}

// GetServicesVarMap compiles the variables passed into 'terraform apply' for services
func GetServicesVarMap(deployConfig deploy.Config, secrets types.Secrets, serviceDockerImagesMap map[string]string) (map[string]string, error) {
	varMap, err := getServicesEnvVarVarMap(deployConfig, secrets)
	if err != nil {
		return map[string]string{}, errors.Wrap(err, "cannot compile service environment variables")
	}
	util.Merge(varMap, getServicesDockerImageVarMap(deployConfig, serviceDockerImagesMap))
	util.Merge(varMap, getServicesURLVarMap(deployConfig))
	util.Merge(varMap, getSharedVarMap(deployConfig, secrets))
	return varMap, nil
}

// getServicesDockerImageVarMap compiles the docker image variables for each service
func getServicesDockerImageVarMap(deployConfig deploy.Config, serviceDockerImagesMap map[string]string) map[string]string {
	varMap := map[string]string{}
	for serviceRole := range deployConfig.AppContext.ServiceContexts {
		varMap[fmt.Sprintf("%s_docker_image", serviceRole)] = serviceDockerImagesMap[serviceRole]
	}
	return varMap
}

// getServiceEnvVarsVarMap compiles env vars for each service
func getServicesEnvVarVarMap(deployConfig deploy.Config, secrets types.Secrets) (map[string]string, error) {
	envVars := map[string]string{}
	serviceEndpoints := endpoints.NewServiceEndpoints(deployConfig.AppContext, types.BuildModeDeploy, deployConfig.RemoteEnvironmentID)
	for serviceRole, serviceContext := range deployConfig.AppContext.ServiceContexts {
		serviceEnvVars := map[string]string{"ROLE": serviceRole}
		util.Merge(serviceEnvVars, deployConfig.AppContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].EnvironmentVariables)
		remoteEnvironment := serviceContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID]
		util.Merge(serviceEnvVars, remoteEnvironment.EnvironmentVariables)
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

func getServicesURLVarMap(deployConfig deploy.Config) map[string]string {
	varMap := map[string]string{}
	for serviceRole, serviceContext := range deployConfig.AppContext.ServiceContexts {
		if serviceContext.Config.Type == types.ServiceTypePublic {
			varMap[fmt.Sprintf("%s_url", serviceRole)] = serviceContext.Config.Remote.Environments[deployConfig.RemoteEnvironmentID].URL
		}
	}
	return varMap
}
