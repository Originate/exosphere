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

func (d *DockerSetup) getDockerEnvVars() map[string]string {
	result := map[string]string{"ROLE": d.Role}
	for _, dependency := range d.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, d.AppConfig, d.AppDir, d.HomeDir)
		for variable, value := range builtDependency.GetServiceEnvVariables() {
			result[variable] = value
		}
	}
	for _, dependency := range d.ServiceConfig.Dependencies {
		result[strings.ToUpper(dependency.Name)] = dependency.Name
	}
	return result
}

func (d *DockerSetup) getDockerLinks() []string {
	result := []string{}
	for _, dependency := range d.ServiceConfig.Dependencies {
		result = append(result, fmt.Sprintf("%s%s:%s", dependency.Name, dependency.Version, dependency.Name))
	}
	return result
}

func (d *DockerSetup) getExternalServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	renderedVolumes, err := dockerHelpers.GetRenderedVolumes(d.ServiceConfig.Docker.Volumes, d.AppConfig.Name, d.Role, d.HomeDir)
	if err != nil {
		return result, err
	}
	result[d.Role] = types.DockerConfig{
		Image:         d.ServiceData.DockerImage,
		ContainerName: d.Role,
		Ports:         d.ServiceConfig.Docker.Ports,
		Environment:   util.JoinStringMaps(d.ServiceConfig.Docker.Environment, d.getDockerEnvVars()),
		Volumes:       renderedVolumes,
		DependsOn:     serviceConfigHelpers.GetServiceDependencies(d.ServiceConfig, d.AppConfig),
	}
	return result, nil
}

func (d *DockerSetup) getInternalServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	result[d.Role] = types.DockerConfig{
		Build:         path.Join("..", d.ServiceData.Location),
		ContainerName: d.Role,
		Command:       d.ServiceConfig.Startup["command"],
		Ports:         d.ServiceConfig.Docker.Ports,
		Links:         d.getDockerLinks(),
		Environment:   d.getDockerEnvVars(),
		DependsOn:     serviceConfigHelpers.GetServiceDependencies(d.ServiceConfig, d.AppConfig),
	}
	dependencyDockerConfigs, err := d.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return util.JoinDockerConfigMaps(result, dependencyDockerConfigs), nil
}

func (d *DockerSetup) getServiceDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range d.ServiceConfig.Dependencies {
		if !dependency.Config.IsEmpty() {
			builtDependency := appDependencyHelpers.Build(dependency, d.AppConfig, d.AppDir, d.HomeDir)
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
func (d *DockerSetup) GetServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerConfigs()
	} else if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerConfigs()
	}
	return map[string]types.DockerConfig{}, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}
