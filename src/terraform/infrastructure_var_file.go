package terraform

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// GenerateInfrastructureVarFile compiles the variable flags passed into a Terraform command
func GenerateInfrastructureVarFile(deployConfig deploy.Config, secrets types.Secrets) error {
	varMap, err := GetInfrastructureVarMap(deployConfig, secrets)
	if err != nil {
		return err
	}
	return writeVarFile(varMap, deployConfig.GetInfrastructureTerraformDir(), deployConfig.RemoteEnvironmentID)
}

// GetInfrastructureVarMap compiles the variables passed into 'terraform apply' for services
func GetInfrastructureVarMap(deployConfig deploy.Config, secrets types.Secrets) (map[string]string, error) {
	varMap, err := getDependenciesEnvVarVarMap(deployConfig)
	if err != nil {
		return map[string]string{}, errors.Wrap(err, "cannot compile dependency variables")
	}
	util.Merge(varMap, getSharedVarMap(deployConfig, secrets))
	return varMap, nil
}

// getDependenciesEnvVarVarMap compiles env vars for each dependency
func getDependenciesEnvVarVarMap(deployConfig deploy.Config) (map[string]string, error) {
	dependencyVars := map[string]string{}
	for dependencyName := range deployConfig.AppContext.GetRemoteDependencies() {
		serviceData, err := getDependencyServiceData(dependencyName, deployConfig)
		if err != nil {
			return map[string]string{}, err
		}
		serviceDataEnvVar, err := createEnvVarString(map[string]string{"SERVICE_DATA": serviceData})
		if err != nil {
			return map[string]string{}, err
		}
		dependencyVars[fmt.Sprintf("%s_env_vars", dependencyName)] = serviceDataEnvVar
	}
	return dependencyVars, nil
}

func getDependencyServiceData(dependencyName string, deployConfig deploy.Config) (string, error) {
	serviceData := deployConfig.AppContext.GetDependencyServiceData(dependencyName)
	serviceDataBytes, err := json.Marshal(serviceData)
	return string(serviceDataBytes), err
}
