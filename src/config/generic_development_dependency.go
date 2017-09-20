package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
)

type genericDevelopmentDependency struct {
	config    types.DevelopmentDependencyConfig
	appConfig types.AppConfig
	appDir    string
	homeDir   string
}

// GetContainerName returns the container name
func (g *genericDevelopmentDependency) GetContainerName() string {
	return g.config.Name + g.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (g *genericDevelopmentDependency) GetDockerConfig() (types.DockerConfig, error) {
	renderedVolumes, err := tools.GetRenderedVolumes(g.config.Config.Volumes, g.appConfig.Name, g.config.Name, g.homeDir)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
		ContainerName: g.GetContainerName(),
		Ports:         g.config.Config.Ports,
		Volumes:       renderedVolumes,
		Environment:   g.config.Config.DependencyEnvironment,
	}, nil
}

// GetOnlineText returns the online text
func (g *genericDevelopmentDependency) GetOnlineText() string {
	return g.config.Config.OnlineText
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *genericDevelopmentDependency) GetServiceEnvVariables() map[string]string {
	return g.config.Config.ServiceEnvironment
}
