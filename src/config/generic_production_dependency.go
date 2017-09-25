package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type genericProductionDependency struct {
	config    types.ProductionDependencyConfig
	appConfig types.AppConfig
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (g *genericProductionDependency) HasDockerConfig() bool {
	return true
}

// GetDockerConfig returns docker configuration and an error if any
func (g *genericProductionDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image: fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
	}, nil
}

func (g *genericProductionDependency) GetServiceName() string {
	return g.config.Name + g.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (g *genericProductionDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": g.config.Version,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (g *genericProductionDependency) GetDeploymentServiceEnvVariables() map[string]string {
	return map[string]string{}
}
