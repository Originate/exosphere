package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type genericDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

//GetDeploymentConfig returns configuration needed in deployment
func (dependency genericDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": dependency.config.Version,
	}
	return config
}
