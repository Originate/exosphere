package config

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

type localGenericDependency struct {
	config types.LocalDependency
}

// GetContainerName returns the container name
func (g *localGenericDependency) GetContainerName() string {
	return g.config.Name + g.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (g *localGenericDependency) GetDockerConfig() (types.DockerConfig, error) {
	volumes := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.config.Name + "_" + path)
		volumes = append(volumes, fmt.Sprintf("%s:%s", name, path))
	}
	return types.DockerConfig{
		Image:       fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
		Ports:       g.config.Config.Ports,
		Volumes:     volumes,
		Environment: g.config.Config.DependencyEnvironment,
		Restart:     "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *localGenericDependency) GetServiceEnvVariables() map[string]string {
	result := map[string]string{}
	result[strings.ToUpper(g.config.Name)] = g.GetContainerName()
	util.Merge(result, g.config.Config.ServiceEnvironment)
	return result
}

// GetVolumeNames returns the named volumes used by this dependency
func (g *localGenericDependency) GetVolumeNames() []string {
	result := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.config.Name + "_" + path)
		result = append(result, name)
	}
	return result
}
