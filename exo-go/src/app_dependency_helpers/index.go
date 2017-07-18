package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type dependency interface {
	getEnvVariables() map[string]string
}

func build(config types.Dependency) dependency {
	switch config.Name {
	case "exocom":
		return ExocomDependency{config}
	}
}
