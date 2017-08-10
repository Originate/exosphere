package config

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/src/docker"
	"github.com/Originate/exosphere/src/types"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
)

func getExternalServiceConfig(serviceDirName string, serviceData types.ServiceData) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	c, err := client.NewEnvClient()
	if err != nil {
		return serviceConfig, err
	}
	yamlFile, err := docker.CatFileInDockerImage(c, serviceData.DockerImage, "service.yml")
	if err != nil {
		return serviceConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the external service '%s'", serviceDirName))
	}
	return serviceConfig, nil
}

func getInternalServiceConfig(appDir, serviceDirName string) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(appDir, serviceDirName, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the internal service '%s'", serviceDirName))
	}
	return serviceConfig, nil
}

// GetInternalServiceConfigs reads the service.yml of all internal services and returns
// the serviceConfig objects and an error (if any)
func GetInternalServiceConfigs(appDir string, appConfig types.AppConfig) (map[string]types.ServiceConfig, error) {
	result := map[string]types.ServiceConfig{}
	for service, serviceData := range GetServiceData(appConfig.Services) {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(appDir, serviceData.Location)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		}
	}
	return result, nil
}

// GetServiceConfigs reads the service.yml of all services and returns
// the serviceConfig objects and an error (if any)
func GetServiceConfigs(appDir string, appConfig types.AppConfig) (map[string]types.ServiceConfig, error) {
	result := map[string]types.ServiceConfig{}
	for service, serviceData := range appConfig.GetServiceData() {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(appDir, serviceData.Location)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		} else if len(serviceData.DockerImage) > 0 {
			serviceConfig, err := getExternalServiceConfig(service, serviceData)
			if err != nil {
				return result, err
			}
			result[service] = serviceConfig
		} else {
			return result, fmt.Errorf("No location or docker image listed for %s", service)
		}
	}
	return result, nil
}

// GetServiceBuiltDependencies returns the AppDependency objects for the
// dependencies defined in the given serviceConfig
func GetServiceBuiltDependencies(serviceConfig types.ServiceConfig, appConfig types.AppConfig, appDir, homeDir string) map[string]AppDependency {
	appBuiltDependencies := GetAppBuiltDependencies(appConfig, appDir, homeDir)
	result := map[string]AppDependency{}
	for _, dependency := range serviceConfig.Dependencies {
		if !dependency.Config.IsEmpty() {
			builtDependency := NewAppDependency(dependency, appConfig, appDir, homeDir)
			result[dependency.Name] = builtDependency
		} else {
			result[dependency.Name] = appBuiltDependencies[dependency.Name]
		}
	}
	return result
}

// GetServiceData returns the configurations data for the given services
func GetServiceData(services types.Services) map[string]types.ServiceData {
	result := make(map[string]types.ServiceData)
	for serviceName, data := range services.Private {
		result[serviceName] = data
	}
	for serviceName, data := range services.Public {
		result[serviceName] = data
	}
	return result
}
