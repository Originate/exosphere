package config

import "github.com/Originate/exosphere/exo-go/src/dockercompose"

// AppDependency contains methods that return config information about a dependency
type AppDependency interface {
	GetContainerName() string
	GetDeploymentConfig() map[string]string
	GetDockerConfig() (dockercompose.DockerConfig, error)
	GetEnvVariables() map[string]string
	GetOnlineText() string
	GetServiceEnvVariables() map[string]string
}

// NewAppDependency returns an AppDependency
func NewAppDependency(dependency Dependency, appConfig AppConfig, appDir, homeDir string) AppDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomDependency{dependency, appConfig, appDir}
	case "nats":
		return &natsDependency{dependency, appConfig, appDir}
	default:
		return &genericDependency{dependency, appConfig, appDir, homeDir}
	}
}
