package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type exocomDevelopmentDependency struct {
	config    types.DevelopmentDependencyConfig
	appConfig types.AppConfig
	appDir    string
}

func (e *exocomDevelopmentDependency) compileServiceRoutes() ([]map[string]interface{}, error) {
	routes := []map[string]interface{}{}
	serviceConfigs, err := GetServiceConfigs(e.appDir, e.appConfig)
	if err != nil {
		return routes, err
	}
	serviceData := e.appConfig.GetServiceData()
	for _, serviceName := range e.appConfig.GetSortedServiceNames() {
		serviceConfig := serviceConfigs[serviceName]
		route := map[string]interface{}{
			"role":     serviceName,
			"receives": serviceConfig.ServiceMessages.Receives,
			"sends":    serviceConfig.ServiceMessages.Sends,
		}
		messageTranslations := serviceData[serviceName].MessageTranslations
		if messageTranslations != nil {
			route["messageTranslations"] = messageTranslations
		}
		routes = append(routes, route)
	}
	return routes, nil
}

// GetContainerName returns the container name
func (e *exocomDevelopmentDependency) GetContainerName() string {
	return e.config.Name + e.config.Version
}

// GetDockerConfig returns docker configuration and an error if any
func (e *exocomDevelopmentDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes, err := e.getServiceRoutesString()
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		ContainerName: e.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Environment: map[string]string{
			"ROLE":           "exocom",
			"SERVICE_ROUTES": serviceRoutes,
		},
	}, nil
}

// GetOnlineText returns the online text for the exocom
func (e *exocomDevelopmentDependency) GetOnlineText() string {
	return "ExoCom online at port"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *exocomDevelopmentDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetContainerName(),
	}
}

func (e *exocomDevelopmentDependency) getServiceRoutesString() (string, error) {
	serviceRoutes, err := e.compileServiceRoutes()
	if err != nil {
		return "", err
	}
	serviceRoutesBytes, err := json.Marshal(serviceRoutes)
	if err != nil {
		return "", err
	}
	return string(serviceRoutesBytes), nil
}
