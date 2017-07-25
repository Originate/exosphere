package appDependencyHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppDependency contains methods that return config information about a dependency
type AppDependency interface {
	GetContainerName() string
	GetDeploymentConfig() map[string]string
	GetDockerConfig() (types.DockerConfig, error)
	GetEnvVariables() map[string]string
	GetOnlineText() string
	GetServiceEnvVariables() map[string]string
}

// Build returns an appDependency
func Build(dependency types.Dependency, appConfig types.AppConfig, appDir, homeDir string) AppDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomDependency{dependency, appConfig, appDir}
	case "nats":
		return &natsDependency{dependency, appConfig, appDir}
	default:
		return &genericDependency{dependency, appConfig, appDir, homeDir}
	}
}
