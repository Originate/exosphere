package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
)

type remoteGenericDependency struct {
	config types.RemoteDependency
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (g *remoteGenericDependency) HasDockerConfig() bool {
	return true
}

// GetDockerConfig returns docker configuration and an error if any
func (g *remoteGenericDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image: fmt.Sprintf("%s:%s", g.config.Name, g.config.Version),
	}, nil
}

func (g *remoteGenericDependency) GetServiceName() string {
	return g.config.Name + g.config.Version
}

//GetDeploymentConfig returns configuration needed in deployment
func (g *remoteGenericDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": g.config.Version,
	}
	return config, nil
}

// GetDeploymentServiceEnvVariables returns configuration needed for each service in deployment
func (g *remoteGenericDependency) GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string {
	return map[string]string{}
}

// GetDeploymentVariables returns a map from string to string of variables that a dependency Terraform module needs
func (g *remoteGenericDependency) GetDeploymentVariables() (map[string]string, error) {
	return map[string]string{}, nil
}
