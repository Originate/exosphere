package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type exocomDependency struct {
	config    types.DependencyConfig
	appConfig types.AppConfig
	appDir    string
}

func (e *exocomDependency) compileServiceRoutes() ([]map[string]interface{}, error) {
	routes := []map[string]interface{}{}
	serviceConfigs, err := GetServiceConfigs(e.appDir, e.appConfig)
	if err != nil {
		return routes, err
	}
	serviceData := e.appConfig.GetServiceData()
	for serviceName, serviceConfig := range serviceConfigs {
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
func (e *exocomDependency) GetContainerName() string {
	return e.config.Name + e.config.Version
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (e *exocomDependency) GetDeploymentConfig() (map[string]string, error) {
	serviceRoutes, err := e.getServiceRoutesString()
	if err != nil {
		return nil, err
	}
	config := map[string]string{
		"version":       e.config.Version,
		"dnsName":       e.appConfig.Production["url"],
		"serviceRoutes": serviceRoutes,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (e *exocomDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": fmt.Sprintf("exocom.%s.local", e.appConfig.Name),
	}
}

// GetDockerConfig returns docker configuration and an error if any
func (e *exocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes, err := e.getServiceRoutesString()
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		ContainerName: e.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Environment: map[string]string{
			"ROLE":           "exocom",
			"PORT":           "80",
			"SERVICE_ROUTES": serviceRoutes,
		},
	}, nil
}

// GetOnlineText returns the online text for the exocom
func (e *exocomDependency) GetOnlineText() string {
	return "ExoCom online at port"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *exocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetContainerName(),
		"EXOCOM_PORT": "80",
	}
}

func (e *exocomDependency) getServiceRoutesString() (string, error) {
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
