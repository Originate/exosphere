package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type postgresDependency struct {
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

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (n *natsDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{}
}

// GetDockerConfig returns docker configuration and an error if any
func (n *natsDependency) GetDockerConfig() (types.DockerConfig, error) {
	ports := []string{"5432:5432"}
	if n.config.Config.Ports > 0 {
		ports = n.config.Config.Ports
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("postgres:%s", n.config.Version),
		ContainerName: n.GetContainerName(),
		Ports:         ports,
	}, nil
}

// GetEnvVariables returns the environment variables
func (n *natsDependency) GetEnvVariables() map[string]string {
	return n.config.Config.DependencyEnvironment
}

// GetOnlineText returns the online text for the nats
func (n *natsDependency) GetOnlineText() string {
	onlineText := "PostgreSQL init process complete"
	if n.config.Config.OnlineText != "" {
		onlineText = n.config.Config.OnlineText
	}
	return onlineText
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (n *natsDependency) GetServiceEnvVariables() map[string]string {
	return n.config.Config.ServiceEnvironment
}
