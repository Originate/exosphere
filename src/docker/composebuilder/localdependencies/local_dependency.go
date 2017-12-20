package localdependencies

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
)

// LocalDependency contains methods that return config information about a local dependency
type LocalDependency struct {
	name       string
	config     types.LocalDependency
	appContext *context.AppContext
}

// NewLocalDependency returns a LocalDependency
func NewLocalDependency(dependencyName string, dependency types.LocalDependency, appContext *context.AppContext) *LocalDependency {
	return &LocalDependency{
		name:       dependencyName,
		config:     dependency,
		appContext: appContext,
	}
}

// GetDockerConfig returns docker configuration and an error if any
func (g *LocalDependency) GetDockerConfig() (types.DockerConfig, error) {
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
		Volumes:     volumes,
		Environment: environment,
		Restart:     "on-failure",
	}, nil
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (g *LocalDependency) GetServiceEnvVariables() map[string]string {
	result := map[string]string{}
	result[fmt.Sprintf("%s_HOST", strings.ToUpper(g.name))] = g.name
	util.Merge(result, g.config.Config.ServiceEnvironment)
	return result
}

// GetVolumeNames returns the named volumes used by this dependency
func (g *LocalDependency) GetVolumeNames() []string {
	result := []string{}
	for _, path := range g.config.Config.Persist {
		name := util.ToSnake(g.name + "_" + path)
		result = append(result, name)
	}
	return result
}

func (g *LocalDependency) getDependencyEnvironment() (map[string]string, error) {
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
