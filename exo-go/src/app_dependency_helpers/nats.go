package appDependencyHelpers

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/types"
)

type natsDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
	appDir    string
}

// GetContainerName returns the container name for the dependency
func (nats natsDependency) GetContainerName() string {
	return nats.config.Name + nats.config.Version
}

// GetDockerConfig returns docker configuration for the dependency and an error if any
func (nats natsDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image:         fmt.Sprintf("nats:%s", nats.config.Version),
		ContainerName: nats.GetContainerName(),
	}, nil
}

// GetEnvVariables returns the environment variables for the depedency
func (nats natsDependency) GetEnvVariables() map[string]string {
	return map[string]string{}
}

// GetOnlineText returns the online text for the nats
func (nats natsDependency) GetOnlineText() string {
	return "Listening for route connections"
}

// GetServiceEnvVariables returns the environment variables for the depedency
func (nats natsDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{"NATS_HOST": nats.GetContainerName()}
}
