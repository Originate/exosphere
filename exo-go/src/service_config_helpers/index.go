package serviceConfigHelpers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/docker/docker/client"
	yaml "gopkg.in/yaml.v2"
)

func GetServiceConfigs(appConfig types.AppConfig) (map[string]types.ServiceConfig, error) {
	serviceConfigs := map[string]types.ServiceConfig{}
	for service, serviceData := range serviceHelpers.GetServiceData(appConfig.Services) {
		if len(serviceData.Location) > 0 {
			serviceConfig, err := getInternalServiceConfig(serviceData.Location)
			if err != nil {
				return serviceConfigs, err
			}
			serviceConfigs[service] = serviceConfig
		} else if len(serviceData.DockerImage) > 0 {
			serviceConfig, err := getExternalServiceConfig(serviceData.DockerImage)
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

func getInternalServiceConfig(serviceLocation string) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(serviceLocation, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		return serviceConfig, err
	}
	return serviceConfig, nil
}

func getExternalServiceConfig(dockerImage string) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	c, err := client.NewEnvClient()
	if err != nil {
		return serviceConfig, err
	}
	yamlFile, err := dockerHelpers.CatFileInDockerImage(c, dockerImage, "service.yml")
	if err != nil {
		return serviceConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		return serviceConfig, err
	}
	return serviceConfig, nil
}
