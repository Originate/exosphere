package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
)

// LocalAppDependency contains methods that return config information about a local dependency
type LocalAppDependency struct {
	name       string
	config     types.LocalDependency
	appContext *context.AppContext
}

// NewLocalAppDependency returns a LocalAppDependency
func NewLocalAppDependency(dependencyName string, dependency types.LocalDependency, appContext *context.AppContext) *LocalAppDependency {
	return &LocalAppDependency{
		name:       dependencyName,
		config:     dependency,
		appContext: appContext,
	}
}

// GetDockerConfig returns docker configuration and an error if any
func (g *LocalAppDependency) GetDockerConfig() (types.DockerConfig, error) {
	volumes := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.name + "_" + path)
		volumes = append(volumes, fmt.Sprintf("%s:%s", name, path))
	}
	environment, err := g.getDependencyEnvironment()
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image:       g.config.Image,
		Ports:       g.config.Config.Ports,
		Volumes:     volumes,
		Environment: environment,
		Restart:     "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *LocalAppDependency) GetServiceEnvVariables() map[string]string {
	result := map[string]string{}
	result[fmt.Sprintf("%s_HOST", strings.ToUpper(g.name))] = g.name
	util.Merge(result, g.config.Config.ServiceEnvironment)
	return result
}

// GetVolumeNames returns the named volumes used by this dependency
func (g *LocalAppDependency) GetVolumeNames() []string {
	result := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.name + "_" + path)
		result = append(result, name)
	}
	return result
}

func (g *LocalAppDependency) getDependencyEnvironment() (map[string]string, error) {
	result := map[string]string{}
	serviceData := g.appContext.GetDependencyServiceData(g.name)
	serviceDataBytes, err := json.Marshal(serviceData)
	if err != nil {
		return result, err
	}
	util.Merge(result, g.config.Config.DependencyEnvironment)
	result["SERVICE_DATA"] = string(serviceDataBytes)
	return result, nil
}
