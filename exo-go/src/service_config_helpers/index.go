package serviceConfigHelpers

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func getExternalServiceConfig(serviceDirName string, serviceData types.ServiceData) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	c, err := client.NewEnvClient()
	if err != nil {
		return serviceConfig, err
	}
	yamlFile, err := dockerHelpers.CatFileInDockerImage(c, serviceData.DockerImage, "service.yml")
	if err != nil {
		return serviceConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the service %s", serviceDirName))
	}
	return serviceConfig, nil
}

func getInternalServiceConfig(serviceDirName string) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(serviceDirName, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return serviceConfig, err
	}
	return serviceConfig, nil
}

// GetServiceData returns a map that maps each service to its serviceData
func getServiceData(services types.Services) map[string]types.ServiceData {
	serviceData := make(map[string]types.ServiceData)
	for service, data := range services.Private {
		serviceData[service] = data
	}
	for service, data := range services.Public {
		serviceData[service] = data
	}
	return serviceData
}

// GetServiceConfigs reads the service.yml of all services and returns
// the serviceConfig objects and an error (if any)
func GetServiceConfigs(appConfig types.AppConfig) (map[string]types.ServiceConfig, error) {
	serviceConfigs := map[string]types.ServiceConfig{}
	for service, serviceData := range getServiceData(appConfig.Services) {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(serviceData.Location)
			if err != nil {
				return serviceConfigs, err
			}
			serviceConfigs[service] = serviceConfig
		} else if len(serviceData.DockerImage) > 0 {
			serviceConfig, err := getExternalServiceConfig(service, serviceData)
			if err != nil {
				return serviceConfigs, err
			}
			serviceConfigs[service] = serviceConfig
		} else {
			return serviceConfigs, fmt.Errorf("No location or docker image listed for %s", service)
		}
	}
	return serviceConfigs, nil
}
