package composebuilder

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// GetBuiltLocalAppDependencies returns the LocalDependency objects for application dependencies only
func GetBuiltLocalAppDependencies(appContext *context.AppContext) map[string]*LocalDependency {
	result := map[string]*LocalDependency{}
	for name, dependency := range appContext.Config.Local.Dependencies {
		builtDependency := NewLocalDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}

// GetBuiltLocalServiceDependencies returns the dependencies for a single service
func GetBuiltLocalServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]*LocalDependency {
	result := map[string]*LocalDependency{}
	for name, dependency := range serviceConfig.Local.Dependencies {
		builtDependency := NewLocalDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}
