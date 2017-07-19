package appDependencyHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppDependency contains methods that return config information about a dependency
type AppDependency interface {
	GetContainerName() string
	GetDockerConfig() (types.DockerConfig, error)
	GetEnvVariables() map[string]string
	GetOnlineText() string
}

// Build returns an appDependency
func Build(dependency types.Dependency, appConfig types.AppConfig) AppDependency {
	switch dependency.Name {
	case "exocom":
		return exocomDependency{dependency, appConfig}
	case "nats":
		return natsDependency{dependency, appConfig}
	default:
		return genericDependency{dependency, appConfig}
	}
}
