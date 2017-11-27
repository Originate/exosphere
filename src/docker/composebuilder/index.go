package composebuilder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
)

// GetDockerComposeProjectName creates a docker compose project name the same way docker-compose mutates the COMPOSE_PROJECT_NAME env var
func GetDockerComposeProjectName(appName string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	replacedStr := reg.ReplaceAllString(appName, "")
	return strings.ToLower(replacedStr)
}

// GetTestDockerComposeProjectName creates a docker compose project name for tests
func GetTestDockerComposeProjectName(appName string) string {
	return GetDockerComposeProjectName(fmt.Sprintf("%stests", appName))
}

// GetApplicationDockerCompose returns the docker compose for a application
func GetApplicationDockerCompose(options ApplicationOptions) (*types.DockerCompose, error) {
	dependencyDockerCompose, err := getDependenciesDockerConfigs(options)
	if err != nil {
		return nil, err
	}
	portReservation := types.NewPortReservation()
	serviceDockerCompose, err := getServicesDockerCompose(options, portReservation)
	if err != nil {
		return nil, err
	}
	return dependencyDockerCompose.Merge(serviceDockerCompose), nil
}

// getDependenciesDockerConfigs returns the docker configs for all the application dependencies
func getDependenciesDockerConfigs(options ApplicationOptions) (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	if options.BuildMode.Type == BuildModeTypeDeploy {
		appDependencies := config.GetBuiltAppProductionDependencies(options.AppConfig, options.AppDir)
		for _, builtDependency := range appDependencies {
			if builtDependency.HasDockerConfig() {
				dockerConfig, err := builtDependency.GetDockerConfig()
				if err != nil {
					return result, err
				}
				result.Services[builtDependency.GetServiceName()] = dockerConfig
			}
		}
	} else {
		appDependencies := config.GetBuiltAppDevelopmentDependencies(options.AppConfig, options.AppDir)
		for _, builtDependency := range appDependencies {
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result.Services[builtDependency.GetContainerName()] = dockerConfig
			for _, name := range builtDependency.GetVolumeNames() {
				result.Volumes[name] = nil
			}
		}
	}
	return result, nil
}

// getServicesDockerCompose returns the docker configs for all the application services
func getServicesDockerCompose(options ApplicationOptions, portReservation *types.PortReservation) (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	serviceConfigs, err := config.GetServiceConfigs(options.AppDir, options.AppConfig)
	if err != nil {
		return result, err
	}
	serviceEndpoints := getServiceEnvVarEndpoints(options, serviceConfigs, portReservation)
	serviceData := options.AppConfig.Services
	for _, serviceRole := range options.AppConfig.GetSortedServiceRoles() {
		serviceDockerCompose, err := GetServiceDockerCompose(options.AppConfig, serviceConfigs[serviceRole], serviceData[serviceRole], serviceRole, options.AppDir, options.BuildMode, serviceEndpoints)
		if err != nil {
			return result, err
		}
		result = result.Merge(serviceDockerCompose)
	}
	return result, nil
}

func getServiceEnvVarEndpoints(options ApplicationOptions, serviceConfigs map[string]types.ServiceConfig, portReservation *types.PortReservation) map[string]*ServiceEndpoints {
	serviceEndpoints := map[string]*ServiceEndpoints{}
	for _, serviceRole := range options.AppConfig.GetSortedServiceRoles() {
		serviceEndpoints[serviceRole] = NewServiceEndpoint(serviceRole, serviceConfigs[serviceRole], portReservation, options.BuildMode)
	}
	return serviceEndpoints
}
