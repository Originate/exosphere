package composebuilder

import (
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// GetServiceDockerConfigs returns the DockerConfigs for a service and its dependencies in docker-compose.yml
func GetServiceDockerConfigs(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, appDir string, homeDir string, mode BuildMode) (types.DockerConfigs, error) {
	if mode.Type == BuildModeTypeDeploy {
		return NewProductionDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir).getServiceDockerConfigs()
	}
	return NewDevelopmentDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir, homeDir, mode).getServiceDockerConfigs()
}

// GetDockerComposeProjectName creates a docker compose project name the same way docker-compose mutates the COMPOSE_PROJECT_NAME env var
func GetDockerComposeProjectName(appName string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	replacedStr := reg.ReplaceAllString(appName, "")
	return strings.ToLower(replacedStr)
}
