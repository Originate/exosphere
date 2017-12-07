package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// LocalAppDependency contains methods that return config information about a dev dependency
type LocalAppDependency interface {
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceEnvVariables() map[string]string
	GetVolumeNames() []string
}

// NewLocalAppDependency returns a LocalAppDependency
func NewLocalAppDependency(dependencyName string, dependency types.LocalDependency, appContext *context.AppContext) LocalAppDependency {
	switch dependency.Type {
	case "exocom":
		return &localExocomDependency{dependencyName, dependency, appContext}
	case "nats":
		return &localNatsDependency{dependencyName, dependency}
	default:
		return &localGenericDependency{dependencyName, dependency}
	}
}
