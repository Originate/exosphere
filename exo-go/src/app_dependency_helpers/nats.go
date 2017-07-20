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

// GetContainerName returns the container name
func (nats natsDependency) GetContainerName() string {
	return nats.config.Name + nats.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (nats natsDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": nats.config.Version,
	}
	return config
}

// GetDockerConfig returns docker configuration and an error if any
func (nats natsDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image:         fmt.Sprintf("nats:%s", nats.config.Version),
		ContainerName: nats.GetContainerName(),
	}, nil
}

// GetEnvVariables returns the environment variables
func (nats natsDependency) GetEnvVariables() map[string]string {
	return map[string]string{}
}

// GetOnlineText returns the online text for the nats
func (nats natsDependency) GetOnlineText() string {
	return "Listening for route connections"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (nats natsDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{"NATS_HOST": nats.GetContainerName()}
}
