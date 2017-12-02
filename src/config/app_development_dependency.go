package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// AppDevelopmentDependency contains methods that return config information about a dev dependency
type AppDevelopmentDependency interface {
	GetContainerName() string
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceEnvVariables() map[string]string
	GetVolumeNames() []string
}

// NewAppDevelopmentDependency returns a AppDevelopmentDependency
func NewAppDevelopmentDependency(dependency types.DevelopmentDependencyConfig, appContext *context.AppContext) AppDevelopmentDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomDevelopmentDependency{dependency, appContext}
	case "nats":
		return &natsDevelopmentDependency{dependency}
	default:
		return &genericDevelopmentDependency{dependency}
	}
}
