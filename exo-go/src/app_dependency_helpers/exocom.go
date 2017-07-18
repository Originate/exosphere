package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type ExocomDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

func (exocom ExocomDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": exocom.config.Version,
		"dnsName": exocom.appConfig.Production["url"],
		//"serviceRoutes":, TODO: wait for exo setup implementation
	}
	return config
}
