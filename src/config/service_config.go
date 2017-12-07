package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// GetBuiltLocalServiceDependencies returns the dependencies for a single service
func GetBuiltLocalServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]LocalAppDependency {
	result := map[string]LocalAppDependency{}
	for name, dependency := range serviceConfig.Local.Dependencies {
		builtDependency := NewLocalAppDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}

// GetBuiltRemoteServiceDependencies returns the dependencies for a single service
func GetBuiltRemoteServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]RemoteAppDependency {
	result := map[string]RemoteAppDependency{}
	for _, dependency := range serviceConfig.Remote.Dependencies {
		builtDependency := NewRemoteAppDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}
