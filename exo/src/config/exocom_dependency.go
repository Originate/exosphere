package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo/src/types"
)

type exocomDependency struct {
	config    types.Dependency
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
func (e *exocomDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": e.config.Version,
		"dnsName": e.appConfig.Production["url"],
		//"serviceRoutes":, TODO: wait for exo setup implementation
		//"dockerImage":, TODO: wait for ecr implementation
	}
	return config
}

// GetDockerConfig returns docker configuration and an error if any
func (e *exocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes, err := e.compileServiceRoutes()
	if err != nil {
		return types.DockerConfig{}, err
	}
	serviceRoutesBytes, err := json.Marshal(serviceRoutes)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		ContainerName: e.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", e.config.Version),
		Command:       "bin/exocom",
		Environment: map[string]string{
			"ROLE":           "exocom",
			"PORT":           "$EXOCOM_PORT",
			"SERVICE_ROUTES": string(serviceRoutesBytes),
		},
	}, nil
}

// GetEnvVariables returns the environment variables
func (e *exocomDependency) GetEnvVariables() map[string]string {
	port := os.Getenv("EXOCOM_PORT")
	if len(port) == 0 {
		port = "80"
	}
	return map[string]string{"EXOCOM_PORT": port}
}

// GetOnlineText returns the online text for the exocom
func (e *exocomDependency) GetOnlineText() string {
	return "ExoCom WebSocket listener online"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (e *exocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": e.GetContainerName(),
		"EXOCOM_PORT": "$EXOCOM_PORT",
	}
}
