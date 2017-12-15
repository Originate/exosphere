package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type remoteExocomDependency struct {
	config     types.RemoteDependency
	appContext *context.AppContext
}

// GetDeploymentConfig returns Exocom configuration needed in deployment
func (e *remoteExocomDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": e.config.Config.Version,
		"dnsName": e.appContext.Config.Remote.URL,
	}
	return config, nil
}
