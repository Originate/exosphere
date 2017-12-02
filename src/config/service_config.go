package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// GetBuiltServiceDevelopmentDependencies returns the dependencies for a single service
func GetBuiltServiceDevelopmentDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]AppDevelopmentDependency {
	result := map[string]AppDevelopmentDependency{}
	for _, dependency := range serviceConfig.Development.Dependencies {
		builtDependency := NewAppDevelopmentDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetBuiltServiceProductionDependencies returns the dependencies for a single service
func GetBuiltServiceProductionDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]AppProductionDependency {
	result := map[string]AppProductionDependency{}
	for _, dependency := range serviceConfig.Production.Dependencies {
		builtDependency := NewAppProductionDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}
