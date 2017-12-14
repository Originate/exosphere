package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// GetBuiltLocalServiceDependencies returns the dependencies for a single service
func GetBuiltLocalServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]*LocalAppDependency {
	result := map[string]*LocalAppDependency{}
	for name, dependency := range serviceConfig.Local.Dependencies {
		builtDependency := NewLocalAppDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}

// GetBuiltRemoteServiceDependencies returns the dependencies for a single service
func GetBuiltRemoteServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext, remoteID string) map[string]RemoteAppDependency {
	result := map[string]RemoteAppDependency{}
	for dependencyName, dependency := range serviceConfig.Remote[remoteID].Dependencies {
		builtDependency := NewRemoteAppDependency(dependencyName, dependency, appContext, remoteID)
		result[dependencyName] = builtDependency
	}
	return result
}
