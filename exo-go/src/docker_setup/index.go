package dockerSetup

import (
	"fmt"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_helpers"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// DockerSetup renders docker-compose.yml file with service configuration
type DockerSetup struct {
	AppConfig     types.AppConfig
	ServiceConfig types.ServiceConfig
	ServiceData   types.ServiceData
	Role          string
	Logger        *logger.Logger
	AppDir        string
	HomeDir       string
}

func (dockerSetup *DockerSetup) getDockerEnvVars() map[string]string {
	result := map[string]string{"ROLE": dockerSetup.Role}
	for _, dependency := range dockerSetup.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, dockerSetup.AppConfig, dockerSetup.AppDir, dockerSetup.HomeDir)
		for variable, value := range builtDependency.GetServiceEnvVariables() {
			result[variable] = value
		}
	}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		result[strings.ToUpper(dependency.Name)] = dependency.Name
	}
	return result
}

func (dockerSetup *DockerSetup) getDockerLinks() []string {
	result := []string{}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		result = append(result, fmt.Sprintf("%s%s:%s", dependency.Name, dependency.Version, dependency.Name))
	}
	return result
}

func (dockerSetup *DockerSetup) getExternalServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	renderedVolumes, err := dockerHelpers.GetRenderedVolumes(dockerSetup.ServiceConfig.Docker.Volumes, dockerSetup.AppConfig.Name, dockerSetup.Role, dockerSetup.HomeDir)
	if err != nil {
		return result, err
	}
	result[dockerSetup.Role] = types.DockerConfig{
		Image:         dockerSetup.ServiceData.DockerImage,
		ContainerName: dockerSetup.Role,
		Ports:         dockerSetup.ServiceConfig.Docker.Ports,
		Environment:   util.JoinStringMaps(dockerSetup.ServiceConfig.Docker.Environment, dockerSetup.getDockerEnvVars()),
		Volumes:       renderedVolumes,
		DependsOn:     serviceConfigHelpers.GetServiceDependencies(dockerSetup.ServiceConfig, dockerSetup.AppConfig),
	}
	return result, nil
}

func (dockerSetup *DockerSetup) getInternalServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	result[dockerSetup.Role] = types.DockerConfig{
		Build:         path.Join("..", dockerSetup.ServiceData.Location),
		ContainerName: dockerSetup.Role,
		Command:       dockerSetup.ServiceConfig.Startup["command"],
		Ports:         dockerSetup.ServiceConfig.Docker.Ports,
		Links:         dockerSetup.getDockerLinks(),
		Environment:   dockerSetup.getDockerEnvVars(),
		DependsOn:     serviceConfigHelpers.GetServiceDependencies(dockerSetup.ServiceConfig, dockerSetup.AppConfig),
	}
	dependencyDockerConfigs, err := dockerSetup.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return util.JoinDockerConfigMaps(result, dependencyDockerConfigs), nil
}

func (dockerSetup *DockerSetup) getServiceDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		if !dependency.Config.IsEmpty() {
			builtDependency := appDependencyHelpers.Build(dependency, dockerSetup.AppConfig, dockerSetup.AppDir, dockerSetup.HomeDir)
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result[builtDependency.GetContainerName()] = dockerConfig
		}
	}
	return result, nil
}

// GetServiceDockerConfigs returns a map the service and its dependencies to their docker configs
func (dockerSetup *DockerSetup) GetServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	if dockerSetup.ServiceData.Location != "" {
		return dockerSetup.getInternalServiceDockerConfigs()
	} else if dockerSetup.ServiceData.DockerImage != "" {
		return dockerSetup.getExternalServiceDockerConfigs()
	}
	return map[string]types.DockerConfig{}, fmt.Errorf("No location or docker image listed for '%s'", dockerSetup.Role)
}
