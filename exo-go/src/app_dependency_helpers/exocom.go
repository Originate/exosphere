package appDependencyHelpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
)

type exocomDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
	appDir    string
}

func (exocom exocomDependency) compileServiceRoutes() ([]map[string]interface{}, error) {
	routes := []map[string]interface{}{}
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(exocom.appDir, exocom.appConfig)
	if err != nil {
		return routes, err
	}
	serviceData := serviceConfigHelpers.GetServiceData(exocom.appConfig.Services)
	for serviceName, serviceConfig := range serviceConfigs {
		route := map[string]interface{}{
			"role":     serviceName,
			"receives": serviceConfig.ServiceMessages.Receives,
			"sends":    serviceConfig.ServiceMessages.Sends,
		}
		namespace := serviceData[serviceName].NameSpace
		if namespace != "" {
			route["namespace"] = namespace
		}
		routes = append(routes, route)
	}
	return routes, nil
}

// GetContainerName returns the container name
func (exocom exocomDependency) GetContainerName() string {
	return exocom.config.Name + exocom.config.Version
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (exocom exocomDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": exocom.config.Version,
		"dnsName": exocom.appConfig.Production["url"],
		//"serviceRoutes":, TODO: wait for exo setup implementation
		//"dockerImage":, TODO: wait for ecr implementation
	}
	return config
}

// GetDockerConfig returns docker configuration and an error if any
func (exocom exocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes, err := exocom.compileServiceRoutes()
	if err != nil {
		return types.DockerConfig{}, err
	}
	serviceRoutesBytes, err := json.Marshal(serviceRoutes)
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		ContainerName: exocom.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", exocom.config.Version),
		Command:       "bin/exocom",
		Environment: map[string]string{
			"ROLE":           "exocom",
			"PORT":           "$EXOCOM_PORT",
			"SERVICE_ROUTES": string(serviceRoutesBytes),
		},
	}, nil
}

// GetEnvVariables returns the environment variables
func (exocom exocomDependency) GetEnvVariables() map[string]string {
	port := os.Getenv("EXOCOM_PORT")
	if len(port) == 0 {
		port = "80"
	}
	return map[string]string{"EXOCOM_PORT": port}
}

// GetOnlineText returns the online text for the exocom
func (exocom exocomDependency) GetOnlineText() string {
	return "ExoCom WebSocket listener online"
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (exocom exocomDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{
		"EXOCOM_HOST": exocom.GetContainerName(),
		"EXOCOM_PORT": "$EXOCOM_PORT",
	}
}
