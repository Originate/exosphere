package composebuilder

import (
	"fmt"
	"os"
	"path"
	"sort"
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
	Mode                     BuildMode
	BuiltAppDependencies     map[string]config.AppDevelopmentDependency
	BuiltServiceDependencies map[string]config.AppDevelopmentDependency
	Role                     string
	AppDir                   string
	HomeDir                  string
	ServiceEndpoints         map[string]*ServiceEndpoints
}

// NewDevelopmentDockerComposeBuilder is DevelopmentDockerComposeBuilder's constructor
func NewDevelopmentDockerComposeBuilder(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role, appDir, homeDir string, mode BuildMode, serviceEndpoints map[string]*ServiceEndpoints) *DevelopmentDockerComposeBuilder {
	return &DevelopmentDockerComposeBuilder{
		AppConfig:                appConfig,
		ServiceConfig:            serviceConfig,
		ServiceData:              serviceData,
		BuiltAppDependencies:     config.GetBuiltAppDevelopmentDependencies(appConfig, appDir, homeDir),
		BuiltServiceDependencies: config.GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir, homeDir),
		Role:             role,
		AppDir:           appDir,
		HomeDir:          homeDir,
		Mode:             mode,
		ServiceEndpoints: serviceEndpoints,
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

func (d *DevelopmentDockerComposeBuilder) getDockerfileName() string {
	if d.Mode.Environment == BuildModeEnvironmentProduction {
		return "Dockerfile.prod"
	}
	return "Dockerfile.dev"
}

func (d *DevelopmentDockerComposeBuilder) getDockerCommand() string {
	switch d.Mode.Environment {
	case BuildModeEnvironmentProduction:
		return ""
	case BuildModeEnvironmentTest:
		return d.ServiceConfig.Development.Scripts["test"]
	default:
		return d.ServiceConfig.Development.Scripts["run"]
	}
}

func (d *DevelopmentDockerComposeBuilder) getDockerPorts() []string {
	switch d.Mode.Environment {
	case BuildModeEnvironmentProduction:
		fallthrough
	case BuildModeEnvironmentDevelopment:
		return d.ServiceEndpoints[d.Role].GetPortMappings()
	default:
		return []string{}
	}
}

func (d *DevelopmentDockerComposeBuilder) getDockerVolumes() []string {
	if !d.Mode.Mount {
		return []string{}
	}
	return []string{d.getServiceFilePath() + ":" + "/mnt"}
}

func (d *DevelopmentDockerComposeBuilder) getRestartPolicy() string {
	if d.Mode.Environment != BuildModeEnvironmentTest {
		return "on-failure"
	}
	return ""
}

func (d *DevelopmentDockerComposeBuilder) getInternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Build: map[string]string{
			"context":    d.getServiceFilePath(),
			"dockerfile": d.getDockerfileName(),
		},
		ContainerName: d.Role,
		Command:       d.getDockerCommand(),
		Ports:         d.getDockerPorts(),
		Links:         d.getDockerLinks(),
		Volumes:       d.getDockerVolumes(),
		Environment:   d.getDockerEnvVars(),
		DependsOn:     d.getServiceDependsOn(),
		Restart:       d.getRestartPolicy(),
	}
	dependencyDockerConfigs, err := d.getServiceDependenciesDockerConfigs()
	if err != nil {
		return result, err
	}
	return result.Merge(dependencyDockerConfigs), nil
}

func (d *DevelopmentDockerComposeBuilder) getExternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	if d.Mode.Environment == BuildModeEnvironmentTest {
		return result, nil
	}
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
		Restart:       d.getRestartPolicy(),
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
	envVars, secrets := d.ServiceConfig.GetEnvVars("development")
	util.Merge(result, envVars)
	for _, secret := range secrets {
		result[secret] = os.Getenv(secret)
	}
	serviceEndpoints := d.createServiceEndpointEnvVars()
	util.Merge(result, serviceEndpoints)
	return result
}

func (d *DevelopmentDockerComposeBuilder) createServiceEndpointEnvVars() map[string]string {
	endpoints := map[string]string{}
	for role, serviceEndpoint := range d.ServiceEndpoints {
		if role == d.Role {
			continue
		}
		util.Merge(endpoints, serviceEndpoint.GetEndpointMappings())
	}
	return endpoints
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
	sort.Strings(result)
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
