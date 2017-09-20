package composebuilder

import (
	"path"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/src/types"
)

// GetServiceDockerConfigs returns the DockerConfigs for a service and its dependencies in docker-compose.yml
func GetServiceDockerConfigs(appConfig types.AppConfig, serviceConfig types.ServiceConfig, serviceData types.ServiceData, role string, appDir string, homeDir string, production bool) (types.DockerConfigs, error) {
	if production {
		return NewProductionDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir).getServiceDockerConfigs()
	}
	return NewDevelopmentDockerComposeBuilder(appConfig, serviceConfig, serviceData, role, appDir, homeDir).getServiceDockerConfigs()
}

// GetDockerComposeProjectName creates a docker compose project name the same way docker-compose mutates the COMPOSE_PROJECT_NAME env var
func GetDockerComposeProjectName(appDir string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	replacedStr := reg.ReplaceAllString(path.Base(appDir), "")
	return strings.ToLower(replacedStr)
}
