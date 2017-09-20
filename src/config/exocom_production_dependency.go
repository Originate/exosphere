package config

import (
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

func (e *exocomProductionDependency) getServiceRoutesString() (string, error) {
	exocomDevelopmentDependency := &exocomDevelopmentDependency{types.DevelopmentDependencyConfig{}, e.appConfig, e.appDir}
	return exocomDevelopmentDependency.getServiceRoutesString()
}
