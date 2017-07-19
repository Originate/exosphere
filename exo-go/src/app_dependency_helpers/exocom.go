package appDependencyHelper

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

type exocomDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (exocom exocomDependency) GetDeploymentConfig() map[string]string {
	config := map[string]string{
		"version": exocom.config.Version,
		"dnsName": exocom.appConfig.Production["url"],
		//"serviceRoutes":, TODO: wait for exo setup implementation
		//"dockerImage":, TODO: wait for ecr implementation
	}
	return config
}
