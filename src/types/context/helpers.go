package context

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func getExternalServiceConfig(dockerImage string) (types.ServiceConfig, error) {
	var serviceConfig types.ServiceConfig
	c, err := client.NewEnvClient()
	if err != nil {
		return serviceConfig, err
	}
	yamlFile, err := tools.CatFileInDockerImage(c, dockerImage, "service.yml")
	if err != nil {
		return serviceConfig, err
	}
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for '%s'", dockerImage))
	}
	return serviceConfig, serviceConfig.ValidateServiceConfig()
}
