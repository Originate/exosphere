package config

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

type genericDevelopmentDependency struct {
	config    types.DevelopmentDependencyConfig
	appConfig types.AppConfig
	appDir    string
}

// GetContainerName returns the container name
func (g *genericDevelopmentDependency) GetContainerName() string {
	return g.config.Name + g.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (g *genericDevelopmentDependency) GetDockerConfig() (types.DockerConfig, error) {
	volumes := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.config.Name + "_" + path)
		volumes = append(volumes, fmt.Sprintf("%s:%s", name, path))
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
		ContainerName: g.GetContainerName(),
		Ports:         g.config.Config.Ports,
		Volumes:       volumes,
		Environment:   g.config.Config.DependencyEnvironment,
		Restart:       "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *genericDevelopmentDependency) GetServiceEnvVariables() map[string]string {
	result := map[string]string{}
	result[strings.ToUpper(g.config.Name)] = g.GetContainerName()
	util.Merge(result, g.config.Config.ServiceEnvironment)
	return result
}

// GetVolumeNames returns the named volumes used by this dependency
func (g *genericDevelopmentDependency) GetVolumeNames() []string {
	result := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.config.Name + "_" + path)
		result = append(result, name)
	}
	return result
}
