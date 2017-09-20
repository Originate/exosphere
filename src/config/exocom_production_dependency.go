package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type exocomProductionDependency struct {
	config    types.ProductionDependencyConfig
	appConfig types.AppConfig
	appDir    string
}

// GetDockerConfig returns docker configuration and an error if any
func (e *exocomProductionDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image: fmt.Sprintf("originate/exocom:%s", e.config.Version),
	}, nil
}

func (e *exocomProductionDependency) GetServiceName() string {
	return e.config.Name + e.config.Version
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (e *exocomProductionDependency) GetDeploymentConfig() (map[string]string, error) {
	serviceRoutes, err := e.getServiceRoutesString()
	if err != nil {
		return nil, err
	}
	config := map[string]string{
		"version":       e.config.Version,
		"dnsName":       e.appConfig.Production.URL,
		"serviceRoutes": serviceRoutes,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (e *exocomProductionDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": fmt.Sprintf("exocom.%s.local", e.appConfig.Name),
	}
}

func (e *exocomProductionDependency) compileServiceRoutes() ([]map[string]interface{}, error) {
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

func (e *exocomProductionDependency) getServiceRoutesString() (string, error) {
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
