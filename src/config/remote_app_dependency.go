package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// RemoteAppDependency contains methods that return config information about a dependency
type RemoteAppDependency interface {
	HasDockerConfig() bool
	GetDockerConfig() (types.DockerConfig, error)
	GetDeploymentConfig() (map[string]string, error)
	GetDeploymentServiceEnvVariables(secrets types.Secrets) map[string]string
	GetDeploymentVariables() (map[string]string, error)
}

// NewRemoteAppDependency returns an AppProductionDependency
func NewRemoteAppDependency(dependencyName string, dependency types.RemoteDependency, appContext *context.AppContext, remoteName string) RemoteAppDependency {
	switch dependency.Type {
	case "exocom":
		return &remoteExocomDependency{dependency, appContext, remoteName}
	case "rds":
		return &remoteRdsDependency{dependencyName, dependency, appContext}
	default:
		return &remoteGenericDependency{dependency}
	}
}
