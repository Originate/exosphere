package config

import (
	"encoding/json"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type localExocomDependency struct {
	name       string
	config     types.LocalDependency
	appContext *context.AppContext
}

// GetDockerConfig returns docker configuration and an error if any
func (e *localExocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceData := e.appContext.GetDependencyServiceData("exocom")
	serviceDataBytes, err := json.Marshal(serviceData)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image: e.config.Image,
		Environment: map[string]string{
			"ROLE":         "exocom",
			"SERVICE_DATA": string(serviceDataBytes),
		},
		Restart: "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *localExocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.name,
	}
}

// GetVolumeNames returns the named volumes used by this dependency
func (e *localExocomDependency) GetVolumeNames() []string {
	return []string{}
}
