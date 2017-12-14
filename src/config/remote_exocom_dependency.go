package config

import (
	"fmt"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type remoteExocomDependency struct {
	config     types.RemoteDependency
	appContext *context.AppContext
}

// HasDockerConfig returns a boolean indicating if a docker-compose.yml entry should be generated for the dependency
func (e *remoteExocomDependency) HasDockerConfig() bool {
	return true
}

// GetDockerConfig returns docker configuration and an error if any
func (e *remoteExocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image: fmt.Sprintf("originate/exocom:%s", e.config.Config.Version),
	}, nil
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (e *remoteExocomDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": e.config.Config.Version,
		"dnsName": e.appContext.Config.Remote.URL,
	}
	return config, nil
}
