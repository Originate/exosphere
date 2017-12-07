package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type localExocomDependency struct {
	config     types.LocalDependency
	appContext *context.AppContext
}

// GetDockerConfig returns docker configuration and an error if any
func (e *localExocomDependency) GetDockerConfig(serviceData map[string]interface{}) (types.DockerConfig, error) {
	serviceDataBytes, err := json.Marshal(serviceData)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image: fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Environment: map[string]string{
			"SERVICE_DATA": string(serviceDataBytes),
		},
		Restart: "on-failure",
	}, nil
}

// GetServiceName returns the service name
func (e *localExocomDependency) GetServiceName() string {
	return e.config.Name + e.config.Version
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *localExocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetServiceName(),
	}
}

// GetVolumeNames returns the named volumes used by this dependency
func (e *localExocomDependency) GetVolumeNames() []string {
	return []string{}
}
