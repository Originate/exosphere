package composebuilder

import (
	"fmt"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/tools"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// DockerComposeBuilder returns the given service config that will appear
// in docker-compose.yml
type DockerComposeBuilder struct {
	AppConfig                types.AppConfig
	ServiceConfig            types.ServiceConfig
	ServiceData              types.ServiceData
	BuiltAppDependencies     map[string]config.AppDependency
	BuiltServiceDependencies map[string]config.AppDependency
	Role                     string
	HomeDir                  string
}

// NewDockerComposeBuilder is DockerComposeBuilder's constructor
func NewDockerComposeBuilder(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, appDir string, homeDir string) *DockerComposeBuilder {
	return &DockerComposeBuilder{
		AppConfig:                appConfig,
		ServiceConfig:            serviceConfig,
		ServiceData:              serviceData,
		BuiltAppDependencies:     config.GetAppBuiltDependencies(appConfig, appDir, homeDir),
		BuiltServiceDependencies: config.GetServiceBuiltDependencies(serviceConfig, appConfig, appDir, homeDir),
		Role:    role,
		HomeDir: homeDir,
	}
}

func (d *DockerComposeBuilder) getDockerEnvVars() map[string]string {
	result := map[string]string{"ROLE": d.Role}
	for _, builtDependency := range d.BuiltAppDependencies {
		for variable, value := range builtDependency.GetServiceEnvVariables() {
			result[variable] = value
		}
	}
	for _, dependency := range d.ServiceConfig.Dependencies {
		result[strings.ToUpper(dependency.Name)] = dependency.Name
	}
	return result
}

func (d *DockerComposeBuilder) getDockerLinks() []string {
	result := []string{}
	for _, dependency := range d.ServiceConfig.Dependencies {
		result = append(result, fmt.Sprintf("%s%s:%s", dependency.Name, dependency.Version, dependency.Name))
	}
	return result
}

func (d *DockerComposeBuilder) getExternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	renderedVolumes, err := tools.GetRenderedVolumes(d.ServiceConfig.Docker.Volumes, d.AppConfig.Name, d.Role, d.HomeDir)
	if err != nil {
		return result, err
	}
	result[d.Role] = types.DockerConfig{
		Image:         d.ServiceData.DockerImage,
		ContainerName: d.Role,
		Ports:         d.ServiceConfig.Docker.Ports,
		Environment:   util.JoinStringMaps(d.ServiceConfig.Docker.Environment, d.getDockerEnvVars()),
		Volumes:       renderedVolumes,
		DependsOn:     d.getServiceDependencyContainerNames(),
	}
	return result, nil
}

func (d *DockerComposeBuilder) getInternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Build:         map[string]string{"context": path.Join("..", d.ServiceData.Location)},
		ContainerName: d.Role,
		Command:       d.ServiceConfig.Startup["command"],
		Ports:         d.ServiceConfig.Docker.Ports,
		Links:         d.getDockerLinks(),
		Environment:   d.getDockerEnvVars(),
		DependsOn:     d.getServiceDependencyContainerNames(),
	}
	dependencyDockerConfigs, err := d.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return result.Merge(dependencyDockerConfigs), nil
}

func (d *DockerComposeBuilder) getServiceDependencyContainerNames() []string {
	result := []string{}
	for _, builtDependency := range d.BuiltAppDependencies {
		result = append(result, builtDependency.GetContainerName())
	}
	for _, builtDependency := range d.BuiltServiceDependencies {
		containerName := builtDependency.GetContainerName()
		if !util.DoesStringArrayContain(result, containerName) {
			result = append(result, containerName)
		}
	}
	return result
}

func (d *DockerComposeBuilder) getServiceDependenciesDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	for _, builtDependency := range d.BuiltServiceDependencies {
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

// GetServiceDockerConfigs returns a map the service and its dependencies to their docker configs
func (d *DockerComposeBuilder) GetServiceDockerConfigs() (types.DockerConfigs, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerConfigs()
	}
	if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerConfigs()
	}
	return types.DockerConfigs{}, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}
