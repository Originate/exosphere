package composebuilder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/endpoints"
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
	serviceDockerCompose, err := getServicesDockerCompose(options)
	if err != nil {
		return nil, err
	}
	return dependencyDockerCompose.Merge(serviceDockerCompose), nil
}

// getDependenciesDockerConfigs returns the docker configs for all the application dependencies
func getDependenciesDockerConfigs(options ApplicationOptions) (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	appDependencies := config.GetBuiltLocalAppDependencies(options.AppContext)
	for dependencyName, builtDependency := range appDependencies {
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result.Services[dependencyName] = dockerConfig
		for _, name := range builtDependency.GetVolumeNames() {
			result.Volumes[name] = nil
		}
	}
	return result, nil
}

// getServicesDockerCompose returns the docker configs for all the application services
func getServicesDockerCompose(options ApplicationOptions) (*types.DockerCompose, error) {
	result := types.NewDockerCompose()
	serviceEndpoints := endpoints.NewServiceEndpoints(options.AppContext, options.BuildMode)
	for _, serviceRole := range options.AppContext.Config.GetSortedServiceRoles() {
		serviceDockerCompose, err := GetServiceDockerCompose(options.AppContext, serviceRole, options.BuildMode, serviceEndpoints)
		if err != nil {
			return result, err
		}
		result = result.Merge(serviceDockerCompose)
	}
	return result, nil
}
