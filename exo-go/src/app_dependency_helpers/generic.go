package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type GenericDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

func (dependency GenericDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": dependency.config.Version,
	}
	return config
}
