package appDependencyHelpers

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
)

type genericDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
	appDir    string
	homeDir   string
}

// GetContainerName returns the container name
func (dependency genericDependency) GetContainerName() string {
	return dependency.config.Name + dependency.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (dependency genericDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": dependency.config.Version,
	}
	return config
}

// GetDockerConfig returns docker configuration and an error if any
func (dependency genericDependency) GetDockerConfig() (types.DockerConfig, error) {
	renderedVolumes, err := dockerHelpers.GetRenderedVolumes(dependency.config.Config.Volumes, dependency.appConfig.Name, dependency.config.Name, dependency.homeDir)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("%s:%s", dependency.config.Name, dependency.config.Version),
		ContainerName: dependency.GetContainerName(),
		Ports:         dependency.config.Config.Ports,
		Volumes:       renderedVolumes,
	}, nil
}

// GetEnvVariables returns the environment variables
func (dependency genericDependency) GetEnvVariables() map[string]string {
	return dependency.config.Config.DependencyEnvironment
}

// GetOnlineText returns the online text
func (dependency genericDependency) GetOnlineText() string {
	return dependency.config.Config.OnlineText
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (dependency genericDependency) GetServiceEnvVariables() map[string]string {
	return dependency.config.Config.ServiceEnvironment
}
