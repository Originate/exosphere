package config

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/docker"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
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

// GetServiceDependencies returns the names of the dependencies of the given serviceConfig
func GetServiceDependencies(serviceConfig types.ServiceConfig, appConfig types.AppConfig) []string {
	result := []string{}
	for _, dependency := range serviceConfig.Dependencies {
		name := dependency.Name + dependency.Version
		if !util.DoesStringArrayContain(result, name) {
			result = append(result, name)
		}
	}
	for _, dependency := range appConfig.Dependencies {
		name := dependency.Name + dependency.Version
		if !util.DoesStringArrayContain(result, name) {
			result = append(result, name)
		}
	}
	return result
}
