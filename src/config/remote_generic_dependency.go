package config

import (
	"github.com/Originate/exosphere/src/types"
)

type remoteGenericDependency struct {
	config types.RemoteDependency
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (g *remoteGenericDependency) HasDockerConfig() bool {
	return false
}

// GetDockerConfig returns docker configuration and an error if any
func (g *remoteGenericDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{}, nil
}

//GetDeploymentConfig returns configuration needed in deployment
func (g *remoteGenericDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": g.config.Config.Version,
	}
	return config, nil
}
