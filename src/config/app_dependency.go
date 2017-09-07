package config

import "github.com/Originate/exosphere/src/types"

// AppDependency contains methods that return config information about a dependency
type AppDependency interface {
	GetContainerName() string
	GetDeploymentConfig() (map[string]string, error)
	GetDeploymentServiceEnvVariables() map[string]string
	GetDockerConfig() (types.DockerConfig, error)
	GetOnlineText() string
	GetServiceEnvVariables() map[string]string
}

// NewAppDependency returns an AppDependency
func NewAppDependency(dependency types.DependencyConfig, appConfig types.AppConfig, appDir, homeDir string) AppDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomDependency{dependency, appConfig, appDir}
	case "nats":
		return &natsDependency{dependency, appConfig, appDir}
	default:
		return &genericDependency{dependency, appConfig, appDir, homeDir}
	}
}
