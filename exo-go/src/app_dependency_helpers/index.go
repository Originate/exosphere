package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type dependency interface {
	GetDeploymentConfig() map[string]string
}

func Build(dependencyConfig types.Dependency, appConfig types.AppConfig) dependency {
	switch dependencyConfig.Name {
	case "exocom":
		return ExocomDependency{dependencyConfig, appConfig}
	default:
		return GenericDependency{dependencyConfig, appConfig}
	}
}
