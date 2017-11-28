package config

import (
	"github.com/Originate/exosphere/src/types"
)

// GetBuiltServiceDevelopmentDependencies returns the dependencies for a single service
func GetBuiltServiceDevelopmentDependencies(serviceConfig types.ServiceConfig, appConfig types.AppConfig, appDir string) map[string]AppDevelopmentDependency {
	result := map[string]AppDevelopmentDependency{}
	for _, dependency := range serviceConfig.Development.Dependencies {
		builtDependency := NewAppDevelopmentDependency(dependency, appConfig, appDir)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetBuiltServiceProductionDependencies returns the dependencies for a single service
func GetBuiltServiceProductionDependencies(serviceConfig types.ServiceConfig, appConfig types.AppConfig, appDir string) map[string]AppProductionDependency {
	result := map[string]AppProductionDependency{}
	for _, dependency := range serviceConfig.Production.Dependencies {
		builtDependency := NewAppProductionDependency(dependency, appConfig, appDir)
		result[dependency.Name] = builtDependency
	}
	return result
}
