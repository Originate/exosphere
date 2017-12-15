package config

import (
	"github.com/Originate/exosphere/src/types"
)

type remoteGenericDependency struct {
	config types.RemoteDependency
}

//GetDeploymentConfig returns configuration needed in deployment
func (g *remoteGenericDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"version": g.config.Config.Version,
	}
	return config, nil
}
