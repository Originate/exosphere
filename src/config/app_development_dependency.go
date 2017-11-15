package config

import "github.com/Originate/exosphere/src/types"

// AppDevelopmentDependency contains methods that return config information about a dev dependency
type AppDevelopmentDependency interface {
	GetContainerName() string
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceEnvVariables() map[string]string
	GetVolumeNames() []string
}

// NewAppDevelopmentDependency returns a AppDevelopmentDependency
func NewAppDevelopmentDependency(dependency types.DevelopmentDependencyConfig, appConfig types.AppConfig, appDir string) AppDevelopmentDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomDevelopmentDependency{dependency, appConfig, appDir}
	case "nats":
		return &natsDevelopmentDependency{dependency, appConfig, appDir}
	default:
		return &genericDevelopmentDependency{dependency, appConfig, appDir}
	}
}
