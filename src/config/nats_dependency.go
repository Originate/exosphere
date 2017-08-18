package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type natsDependency struct {
	config    types.DependencyConfig
	appConfig types.AppConfig
	appDir    string
}

// GetContainerName returns the container name
func (n *natsDependency) GetContainerName() string {
	return n.config.Name + n.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (n *natsDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": n.config.Version,
	}
	return config, nil
}

// GetDockerConfig returns docker configuration and an error if any
func (n *natsDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image:         fmt.Sprintf("nats:%s", n.config.Version),
		ContainerName: n.GetContainerName(),
	}, nil
}

// GetEnvVariables returns the environment variables
func (n *natsDependency) GetEnvVariables() map[string]string {
	return map[string]string{}
}

// GetOnlineText returns the online text for the nats
func (n *natsDependency) GetOnlineText() string {
	return "Listening for route connections"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (n *natsDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{"NATS_HOST": n.GetContainerName()}
}
