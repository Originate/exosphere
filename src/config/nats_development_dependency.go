package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type natsDevelopmentDependency struct {
	config    types.DevelopmentDependencyConfig
	appConfig types.AppConfig
	appDir    string
}

// GetContainerName returns the container name
func (n *natsDevelopmentDependency) GetContainerName() string {
	return n.config.Name + n.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (n *natsDevelopmentDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image:         fmt.Sprintf("nats:%s", n.config.Version),
		ContainerName: n.GetContainerName(),
	}, nil
}

// GetEnvVariables returns the environment variables
func (n *natsDevelopmentDependency) GetEnvVariables() map[string]string {
	return map[string]string{}
}

// GetOnlineText returns the online text for the nats
func (n *natsDevelopmentDependency) GetOnlineText() string {
	return "Listening for route connections"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (n *natsDevelopmentDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{"NATS_HOST": n.GetContainerName()}
}
