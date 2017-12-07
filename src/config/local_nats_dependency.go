package config

import (
	"github.com/Originate/exosphere/src/types"
)

type localNatsDependency struct {
	name   string
	config types.LocalDependency
}

// GetDockerConfig returns docker configuration and an error if any
func (n *localNatsDependency) GetDockerConfig() (types.DockerConfig, error) {
	return types.DockerConfig{
		Image:   n.config.Image,
		Restart: "on-failure",
	}, nil
}

// GetEnvVariables returns the environment variables
func (n *localNatsDependency) GetEnvVariables() map[string]string {
	return map[string]string{}
}

// GetServiceEnvVariables returns the environment variables that need to
// be passed to services that use it
func (n *localNatsDependency) GetServiceEnvVariables() map[string]string {
	return map[string]string{"NATS_HOST": n.name}
}

// GetVolumeNames returns the named volumes used by this dependency
func (n *localNatsDependency) GetVolumeNames() []string {
	return []string{}
}
