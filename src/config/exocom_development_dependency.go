package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type exocomDevelopmentDependency struct {
	config    types.DevelopmentDependencyConfig
	appConfig types.AppConfig
	appDir    string
}

// GetContainerName returns the container name
func (e *exocomDevelopmentDependency) GetContainerName() string {
	return e.config.Name + e.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (e *exocomDevelopmentDependency) GetDockerConfig(serviceData map[string]interface{}) (types.DockerConfig, error) {
	serviceDataBytes, err := json.Marshal(serviceData)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		ContainerName: e.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Environment: map[string]string{
			"SERVICE_DATA": string(serviceDataBytes),
		},
		Restart: "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *exocomDevelopmentDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetContainerName(),
	}
}

// GetVolumeNames returns the named volumes used by this dependency
func (e *exocomDevelopmentDependency) GetVolumeNames() []string {
	return []string{}
}
