package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppDependency contains methods that return config information about a dependency
type AppDependency interface {
	GetDeploymentConfig() map[string]string
}

// Build returns an appDependency
func Build(dependency types.Dependency, appConfig types.AppConfig) AppDependency {
	switch dependency.Name {
	case "exocom":
		return exocomDependency{dependency, appConfig}
	default:
		return genericDependency{dependency, appConfig}
	}
}
