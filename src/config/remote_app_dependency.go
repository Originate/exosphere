package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// RemoteAppDependency contains methods that return config information about a dependency
type RemoteAppDependency interface {
	HasDockerConfig() bool
	GetDockerConfig() (types.DockerConfig, error)
	GetServiceName() string
	GetDeploymentConfig() (map[string]string, error)
	GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string
	GetDeploymentVariables() (map[string]string, error)
}

// NewRemoteAppDependency returns an AppProductionDependency
func NewRemoteAppDependency(dependency types.RemoteDependency, appContext *context.AppContext) RemoteAppDependency {
	switch dependency.Name {
	case "exocom":
		return &remoteExocomDependency{dependency, appContext}
	case "postgres":
		fallthrough
	case "mysql":
		return &remoteRdsDependency{dependency, appContext}
	default:
		return &remoteGenericDependency{dependency}
	}
}
