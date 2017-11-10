package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/types"
	prompt "github.com/kofalt/go-prompt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetBuiltDevelopmentDependencies returns the AppDevelopmentDependency objects for application and service
// dev dependencies of the entire application
func GetBuiltDevelopmentDependencies(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir string) map[string]AppDevelopmentDependency {
	result := GetBuiltAppDevelopmentDependencies(appConfig, appDir)
	for _, serviceConfig := range serviceConfigs {
		for dependencyName, builtDependency := range GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetBuiltAppDevelopmentDependencies returns the AppDevelopmentDependency objects for application dependencies only
func GetBuiltAppDevelopmentDependencies(appConfig types.AppConfig, appDir string) map[string]AppDevelopmentDependency {
	result := map[string]AppDevelopmentDependency{}
	for _, dependency := range appConfig.Development.Dependencies {
		builtDependency := NewAppDevelopmentDependency(dependency, appConfig, appDir)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetBuiltProductionDependencies returns the AppProductionDependency objects for the application and service
// prod dependencies of the entire application
func GetBuiltProductionDependencies(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir string) map[string]AppProductionDependency {
	result := GetBuiltAppProductionDependencies(appConfig, appDir)
	for _, serviceConfig := range serviceConfigs {
		for dependencyName, builtDependency := range GetBuiltServiceProductionDependencies(serviceConfig, appConfig, appDir) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetBuiltAppProductionDependencies returns the AppProductionDependency objects for the application dependencies only
func GetBuiltAppProductionDependencies(appConfig types.AppConfig, appDir string) map[string]AppProductionDependency {
	result := map[string]AppProductionDependency{}
	for _, dependency := range appConfig.Production.Dependencies {
		builtDependency := NewAppProductionDependency(dependency, appConfig, appDir)
		result[dependency.Name] = builtDependency
	}
	return result
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(appDir string, serviceRole string, appConfig types.AppConfig) error {
	protectionLevels := []string{"public", "private", "worker"}
	switch protectionLevels[prompt.Choose("Protection Level", protectionLevels)] {
	case "public":
		if appConfig.Services.Public == nil {
			appConfig.Services.Public = make(map[string]types.ServiceData)
		}
		appConfig.Services.Public[serviceRole] = types.ServiceData{Location: fmt.Sprintf("./%s", serviceRole)}
	case "private":
		if appConfig.Services.Private == nil {
			appConfig.Services.Private = make(map[string]types.ServiceData)
		}
		appConfig.Services.Private[serviceRole] = types.ServiceData{Location: fmt.Sprintf("./%s", serviceRole)}
	case "worker":
		if appConfig.Services.Worker == nil {
			appConfig.Services.Worker = make(map[string]types.ServiceData)
		}
		appConfig.Services.Worker[serviceRole] = types.ServiceData{Location: fmt.Sprintf("./%s", serviceRole)}
	}
	bytes, err := yaml.Marshal(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal application.yml")
	}
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), bytes, 0777)
}
