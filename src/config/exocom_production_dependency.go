package config

import (
	"encoding/json"
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type exocomProductionDependency struct {
	config     types.ProductionDependencyConfig
	appContext *types.AppContext
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
		"dnsName": e.appContext.Config.Production.URL,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (e *exocomProductionDependency) GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string {
	return map[string]string{
		"EXOCOM_HOST": fmt.Sprintf("exocom.%s.local", e.appContext.Config.Name),
	}
}

// GetDeploymentVariables returns a map from string to string of variables that a dependency Terraform module needs
func (e *exocomProductionDependency) GetDeploymentVariables() (map[string]string, error) {
	exocomDevelopmentDependency := &exocomDevelopmentDependency{types.DevelopmentDependencyConfig{}, e.appContext}
	serviceRoutes, err := exocomDevelopmentDependency.compileServiceRoutes()
	if err != nil {
		return map[string]string{}, err
	}
	// marshal the serviceRoutes map[string]interface{} object into string, so that it can be passed into terraform/command_helpers.go/createEnvVarString as a map[string]string
	// this is the proper number of encodings so that the final encoding can be pass as a cli flag to terraform commands
	serviceRoutesJSON, err := json.Marshal(serviceRoutes)
	if err != nil {
		return map[string]string{}, err
	}
	return map[string]string{"SERVICE_ROUTES": string(serviceRoutesJSON)}, err
}
