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

// DevelopmentDockerComposeBuilder contains the docker-compose.yml config for a single service
type DevelopmentDockerComposeBuilder struct {
	AppConfig                types.AppConfig
	ServiceConfig            types.ServiceConfig
	ServiceData              types.ServiceData
	Production               bool
	BuiltAppDependencies     map[string]config.AppDevelopmentDependency
	BuiltServiceDependencies map[string]config.AppDevelopmentDependency
	Role                     string
	AppDir                   string
	HomeDir                  string
}

// NewDevelopmentDockerComposeBuilder is DevelopmentDockerComposeBuilder's constructor
func NewDevelopmentDockerComposeBuilder(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role, appDir, homeDir string) *DevelopmentDockerComposeBuilder {
	return &DevelopmentDockerComposeBuilder{
		AppConfig:                appConfig,
		ServiceConfig:            serviceConfig,
		ServiceData:              serviceData,
		BuiltAppDependencies:     config.GetBuiltAppDevelopmentDependencies(appConfig, appDir, homeDir),
		BuiltServiceDependencies: config.GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir, homeDir),
		Role:    role,
		AppDir:  appDir,
		HomeDir: homeDir,
	}
}

// getServiceDockerConfigs returns a DockerConfig object for a single service and its dependencies (if any(
func (d *DevelopmentDockerComposeBuilder) getServiceDockerConfigs() (types.DockerConfigs, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerConfigs()
	}
	if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerConfigs()
	}
	return types.DockerConfigs{}, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}

func (d *DevelopmentDockerComposeBuilder) getInternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Build: map[string]string{
			"context":    d.getServiceFilePath(),
			"dockerfile": "Dockerfile.dev",
		},
		ContainerName: d.Role,
		Command:       d.ServiceConfig.Development.Scripts["run"],
		Ports:         d.ServiceConfig.Docker.Ports,
		Links:         d.getDockerLinks(),
		Volumes:       []string{d.getServiceFilePath() + ":" + "/mnt"},
		Environment:   d.getDockerEnvVars(),
		DependsOn:     d.getServiceDependsOn(),
	}
	dependencyDockerConfigs, err := d.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return result.Merge(dependencyDockerConfigs), nil
}

func (d *DevelopmentDockerComposeBuilder) getExternalServiceDockerConfigs() (types.DockerConfigs, error) {
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
		DependsOn:     d.getServiceDependsOn(),
	}
	return result, nil
}

func (d *DevelopmentDockerComposeBuilder) getServiceFilePath() string {
	return path.Join(d.AppDir, d.ServiceData.Location)
}

func (d *DevelopmentDockerComposeBuilder) getDockerLinks() []string {
	result := []string{}
	for _, dependency := range d.ServiceConfig.Development.Dependencies {
		result = append(result, fmt.Sprintf("%s%s:%s", dependency.Name, dependency.Version, dependency.Name))
	}
	return result
}

func (d *DevelopmentDockerComposeBuilder) getDockerEnvVars() map[string]string {
	result := map[string]string{"ROLE": d.Role}
	for _, builtDependency := range d.BuiltAppDependencies {
		for variable, value := range builtDependency.GetServiceEnvVariables() {
			result[variable] = value
		}
	}
	for _, builtDependency := range d.BuiltServiceDependencies {
		for variable, value := range builtDependency.GetServiceEnvVariables() {
			result[variable] = value
		}
	}
	for _, dependency := range d.ServiceConfig.Development.Dependencies {
		result[strings.ToUpper(dependency.Name)] = dependency.Name
	}
	return result
}

func (d *DevelopmentDockerComposeBuilder) getServiceDependsOn() []string {
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

// returns the DockerConfigs object for a service's dependencies
func (d *DevelopmentDockerComposeBuilder) getServiceDependenciesDockerConfigs() (types.DockerConfigs, error) {
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
