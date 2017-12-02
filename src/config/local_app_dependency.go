package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// LocalAppDependency contains methods that return config information about a dev dependency
type LocalAppDependency interface {
	GetContainerName() string
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceEnvVariables() map[string]string
	GetVolumeNames() []string
}

// NewLocalAppDependency returns a LocalAppDependency
func NewLocalAppDependency(dependency types.LocalDependency, appContext *context.AppContext) LocalAppDependency {
	switch dependency.Name {
	case "exocom":
		return &localExocomDependency{dependency, appContext}
	case "nats":
		return &localNatsDependency{dependency}
	default:
		return &localGenericDependency{dependency}
	}
}
