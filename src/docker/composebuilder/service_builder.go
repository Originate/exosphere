package composebuilder

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// ServiceComposeBuilder contains the docker-compose.yml config for a single service
type ServiceComposeBuilder struct {
	AppConfig                types.AppConfig
	ServiceConfig            types.ServiceConfig
	ServiceData              types.ServiceData
	Mode                     BuildMode
	BuiltAppDependencies     map[string]config.AppDevelopmentDependency
	BuiltServiceDependencies map[string]config.AppDevelopmentDependency
	Role                     string
	AppDir                   string
	ServiceEndpoints         map[string]*ServiceEndpoints
}

// GetServiceDockerConfig returns the DockerConfigs for a service and its dependencies in docker-compose.yml
func GetServiceDockerConfig(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, appDir string, mode BuildMode, serviceEndpoints map[string]*ServiceEndpoints) (types.DockerConfigs, error) {
	return NewServiceComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir, mode, serviceEndpoints).getServiceDockerConfigs()
}

// NewServiceComposeBuilder is ServiceComposeBuilder's constructor
func NewServiceComposeBuilder(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role, appDir string, mode BuildMode, serviceEndpoints map[string]*ServiceEndpoints) *ServiceComposeBuilder {
	if mode.Environment == BuildModeEnvironmentTest {
		role = appConfig.GetTestRole(role)
	}
	return &ServiceComposeBuilder{
		AppConfig:                appConfig,
		ServiceConfig:            serviceConfig,
		ServiceData:              serviceData,
		BuiltAppDependencies:     config.GetBuiltAppDevelopmentDependencies(appConfig, appDir),
		BuiltServiceDependencies: config.GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir),
		Role:             role,
		AppDir:           appDir,
		Mode:             mode,
		ServiceEndpoints: serviceEndpoints,
	}
}

// getServiceDockerConfigs returns a DockerConfig object for a single service and its dependencies (if any(
func (d *ServiceComposeBuilder) getServiceDockerConfigs() (types.DockerConfigs, error) {
	if d.ServiceData.Location != "" {
		return d.getInternalServiceDockerConfigs()
	}
	if d.ServiceData.DockerImage != "" {
		return d.getExternalServiceDockerConfigs()
	}
	return types.DockerConfigs{}, fmt.Errorf("No location or docker image listed for '%s'", d.Role)
}

func (d *ServiceComposeBuilder) getDockerfileName() string {
	if d.Mode.Environment == BuildModeEnvironmentProduction {
		return "Dockerfile.prod"
	}
	return "Dockerfile.dev"
}

func (d *ServiceComposeBuilder) getDockerCommand() string {
	switch d.Mode.Environment {
	case BuildModeEnvironmentProduction:
		return ""
	case BuildModeEnvironmentTest:
		return d.ServiceConfig.Development.Scripts["test"]
	default:
		return d.ServiceConfig.Development.Scripts["run"]
	}
}

func (d *ServiceComposeBuilder) getDockerPorts() []string {
	switch d.Mode.Environment {
	case BuildModeEnvironmentProduction:
		fallthrough
	case BuildModeEnvironmentDevelopment:
		return d.ServiceEndpoints[d.Role].GetPortMappings()
	default:
		return []string{}
	}
}

func (d *ServiceComposeBuilder) getDockerVolumes() []string {
	if !d.Mode.Mount {
		return []string{}
	}
	return []string{d.getServiceFilePath() + ":" + "/mnt"}
}

func (d *ServiceComposeBuilder) getRestartPolicy() string {
	if d.Mode.Environment != BuildModeEnvironmentTest {
		return "on-failure"
	}
	return ""
}

func (d *ServiceComposeBuilder) getInternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	result[d.Role] = types.DockerConfig{
		Build: map[string]string{
			"context":    d.getServiceFilePath(),
			"dockerfile": d.getDockerfileName(),
		},
		ContainerName: d.Role,
		Command:       d.getDockerCommand(),
		Ports:         d.getDockerPorts(),
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

func (d *ServiceComposeBuilder) getExternalServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	if d.Mode.Environment == BuildModeEnvironmentTest {
		return result, nil
	}
	result[d.Role] = types.DockerConfig{
		Image:         d.ServiceData.DockerImage,
		ContainerName: d.Role,
		Command:       d.getDockerCommand(),
		Ports:         d.getDockerPorts(),
		Environment:   d.getDockerEnvVars(),
		DependsOn:     d.getServiceDependsOn(),
		Restart:       d.getRestartPolicy(),
	}
	return result, nil
}

func (d *ServiceComposeBuilder) getServiceFilePath() string {
	return path.Join("${APP_PATH}", d.ServiceData.Location)
}

func (d *ServiceComposeBuilder) getDockerEnvVars() map[string]string {
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
	envVars, secrets := d.ServiceConfig.GetEnvVars("development")
	util.Merge(result, envVars)
	for _, secret := range secrets {
		result[secret] = os.Getenv(secret)
	}
	serviceEndpoints := d.createServiceEndpointEnvVars()
	util.Merge(result, serviceEndpoints)
	return result
}

func (d *ServiceComposeBuilder) createServiceEndpointEnvVars() map[string]string {
	endpoints := map[string]string{}
	for role, serviceEndpoint := range d.ServiceEndpoints {
		if role == d.Role {
			continue
		}
		util.Merge(endpoints, serviceEndpoint.GetEndpointMappings())
	}
	return endpoints
}

func (d *ServiceComposeBuilder) getServiceDependsOn() []string {
	result := []string{}
	for _, builtDependency := range d.BuiltAppDependencies {
		result = append(result, builtDependency.GetContainerName())
	}
	for _, builtDependency := range d.BuiltServiceDependencies {
		result = append(result, builtDependency.GetContainerName())
	}
	sort.Strings(result)
	return result
}

// returns the DockerConfigs object for a service's dependencies
func (d *ServiceComposeBuilder) getServiceDependenciesDockerConfigs() (types.DockerConfigs, error) {
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
