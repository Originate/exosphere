package config

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func getExternalServiceConfig(serviceDirName string, serviceData ServiceData) (ServiceConfig, error) {
	var serviceConfig ServiceConfig
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
