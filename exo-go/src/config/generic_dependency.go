package config

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker"
	"github.com/Originate/exosphere/exo-go/src/types"
)

type genericDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
	appDir    string
	homeDir   string
}

// GetContainerName returns the container name
func (g *genericDependency) GetContainerName() string {
	return g.config.Name + g.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (g *genericDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": g.config.Version,
	}
	return config
}

// GetDockerConfig returns docker configuration and an error if any
func (g *genericDependency) GetDockerConfig() (types.DockerConfig, error) {
	renderedVolumes, err := docker.GetRenderedVolumes(g.config.Config.Volumes, g.appConfig.Name, g.config.Name, g.homeDir)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
		ContainerName: g.GetContainerName(),
		Ports:         g.config.Config.Ports,
		Volumes:       renderedVolumes,
	}, nil
}

// GetEnvVariables returns the environment variables
func (g *genericDependency) GetEnvVariables() map[string]string {
	return g.config.Config.DependencyEnvironment
}

// GetOnlineText returns the online text
func (g *genericDependency) GetOnlineText() string {
	return g.config.Config.OnlineText
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *genericDependency) GetServiceEnvVariables() map[string]string {
	return g.config.Config.ServiceEnvironment
}
