package composebuilder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// GetServiceDockerCompose returns the DockerCompose for a service and its dependencies in docker-compose.yml
func GetServiceDockerCompose(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, appDir string, mode BuildMode, serviceEndpoints map[string]*ServiceEndpoints) (*types.DockerCompose, error) {
	switch {
	case mode.Type == BuildModeTypeDeploy:
		return NewProductionDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir).getServiceDockerCompose()
	case mode.Environment == BuildModeEnvironmentTest:
		testRole := appConfig.GetTestRole(role)
		return NewDevelopmentDockerComposeBuilder(appConfig, serviceConfig, serviceData, testRole, appDir, mode, serviceEndpoints).getServiceDockerCompose()
	default:
		return NewDevelopmentDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir, mode, serviceEndpoints).getServiceDockerCompose()
	}
}

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
