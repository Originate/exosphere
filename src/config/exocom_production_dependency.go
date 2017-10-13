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

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (e *exocomProductionDependency) HasDockerConfig() bool {
	return true
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
	config := map[string]string{
		"version": e.config.Version,
		"dnsName": e.appConfig.Production.URL,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (e *exocomProductionDependency) GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string {
	return map[string]string{
		"EXOCOM_HOST": fmt.Sprintf("exocom.%s.local", e.appConfig.Name),
	}
}

// GetDeploymentVariables returns a map from string to string of variables that a dependency Terraform module needs
func (e *exocomProductionDependency) GetDeploymentVariables() (map[string]string, error) {
	exocomDevelopmentDependency := &exocomDevelopmentDependency{types.DevelopmentDependencyConfig{}, e.appConfig, e.appDir}
	serviceRoutes, err := exocomDevelopmentDependency.compileServiceRoutes()
	if err != nil {
		return map[string]string{}, err
	}
	serviceRoutesJSON, err := json.Marshal(serviceRoutes)
	if err != nil {
		return map[string]string{}, err
	}
	return map[string]string{"SERVICE_ROUTES": string(serviceRoutesJSON)}, err
}
