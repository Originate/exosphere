package dockerSetup

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/pkg/errors"
)

// DockerSetup renders docker-compose.yml file with service configuration
type DockerSetup struct {
	AppConfig     types.AppConfig
	ServiceConfig types.ServiceConfig
	ServiceData   types.ServiceData
	Role          string
	Logger        *logger.Logger
}

// NewDockerSetup is DockerSetup's constructor
func NewDockerSetup(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, logger *logger.Logger) *DockerSetup {
	return &DockerSetup{
		AppConfig:     appConfig,
		ServiceConfig: serviceConfig,
		Role:          role,
		Logger:        logger,
	}
}

// GetServiceDockerConfig returns a map each service to its doker config
func (dockerSetup *DockerSetup) GetServiceDockerConfig() (map[string]types.DockerConfig, error) {
	if len(dockerSetup.ServiceData.Location) > 0 {
		return dockerSetup.getInternalServiceDockerConfigs()
	} else if len(dockerSetup.ServiceData.DockerImage) > 0 {
		return dockerSetup.getExternalServiceDockerConfigs()
	}
	return map[string]types.DockerConfig{}, fmt.Errorf("No location or docker image listed for '%s'", dockerSetup.Role)
}

func (dockerSetup *DockerSetup) getDockerLinks() []string {
	result := []string{}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		result = append(result, fmt.Sprintf("%s%s:%s", dependency.Name, dependency.Version, dependency.Name))
	}
	return result
}

func (dockerSetup *DockerSetup) getDockerEnvVars() map[string]string {
	result := map[string]string{"ROLE": dockerSetup.Role}
	for _, dependency := range dockerSetup.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, dockerSetup.AppConfig)
		for variable, value := range builtDependency.GetEnvVariables() {
			result[variable] = value
		}
	}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		result[strings.ToUpper(dependency.Name)] = dependency.Name
	}
	return result
}

func (dockerSetup *DockerSetup) getExternalServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	renderedVolumes, err := dockerSetup.getRenderedVolumes(dockerSetup.ServiceConfig.Docker.Volumes)
	if err != nil {
		return result, err
	}
	result[dockerSetup.Role] = types.DockerConfig{
		Image:         dockerSetup.ServiceData.DockerImage,
		ContainerName: dockerSetup.Role,
		Ports:         dockerSetup.ServiceConfig.Docker.Ports,
		Environment:   util.JoinMaps(dockerSetup.ServiceConfig.Docker.Environment, dockerSetup.getDockerEnvVars()),
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
	return joinDockerConfigMaps(result, dependencyDockerConfigs), nil
}

func (dockerSetup *DockerSetup) getRenderedVolumes(volumes []string) ([]string, error) {
	homeDir, err := osHelpers.GetUserHomeDir()
	if err != nil {
		return []string{}, err
	}
	dataPath := path.Join(homeDir, ".exosphere", dockerSetup.AppConfig.Name, dockerSetup.Role, "data")
	renderedVolumes := []string{}
	if err := os.MkdirAll(dataPath, 0777); err != nil { //nolint gas
		return renderedVolumes, errors.Wrap(err, "Failed to create the necessary directories for the volumes")
	}
	for _, volume := range volumes {
		renderedVolumes = append(renderedVolumes, strings.Replace(volume, "{{EXO_DATA_PATH}}", dataPath, -1))
	}
	return renderedVolumes, nil
}

func (dockerSetup *DockerSetup) getServiceDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range dockerSetup.ServiceConfig.Dependencies {
		if !dependency.Config.IsEmpty() {
			builtDependency := appDependencyHelpers.Build(dependency, dockerSetup.AppConfig)
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result[builtDependency.GetContainerName()] = dockerConfig
		}
	}
	return result, nil
}
