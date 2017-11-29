package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// AppProductionDependency contains methods that return config information about a dependency
type AppProductionDependency interface {
	HasDockerConfig() bool
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceName() string
	GetDeploymentConfig() (map[string]string, error)
	GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string
	GetDeploymentVariables() (map[string]string, error)
}

// NewAppProductionDependency returns an AppProductionDependency
func NewAppProductionDependency(dependency types.ProductionDependencyConfig, appContext *context.AppContext) AppProductionDependency {
	switch dependency.Name {
	case "exocom":
		return &exocomProductionDependency{dependency, appContext}
	case "postgres":
		fallthrough
	case "mysql":
		return &rdsProductionDependency{dependency, appContext}
	default:
		return &genericProductionDependency{dependency}
	}
}
