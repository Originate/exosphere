package config

import (
	"github.com/Originate/exosphere/src/types"
)

// AppProductionDependency contains methods that return config information about a dependency
type AppProductionDependency interface {
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceName() string
	GetDeploymentConfig() (map[string]string, error)
	GetDeploymentServiceEnvVariables() map[string]string
}

// NewAppProductionDependency returns an AppProductionDependency
func NewAppProductionDependency(dependency types.ProductionDependencyConfig, appConfig types.AppConfig, appDir string) AppProductionDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomProductionDependency{dependency, appConfig, appDir}
	case "postgres":
		fallthrough
	case "mysql":
		return &rdsProductionDependency{dependency, appConfig}
	default:
		return &genericProductionDependency{dependency, appConfig}
	}
}
