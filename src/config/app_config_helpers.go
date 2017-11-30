package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetBuiltAppDevelopmentDependencies returns the AppDevelopmentDependency objects for application dependencies only
func GetBuiltAppDevelopmentDependencies(appContext *context.AppContext) map[string]AppDevelopmentDependency {
	result := map[string]AppDevelopmentDependency{}
	for _, dependency := range appContext.Config.Development.Dependencies {
		builtDependency := NewAppDevelopmentDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetBuiltProductionDependencies returns the AppProductionDependency objects for the application and service
// prod dependencies of the entire application
func GetBuiltProductionDependencies(appContext *context.AppContext) map[string]AppProductionDependency {
	result := GetBuiltAppProductionDependencies(appContext)
	for _, serviceContext := range appContext.ServiceContexts {
		for dependencyName, builtDependency := range GetBuiltServiceProductionDependencies(serviceContext.Config, appContext) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetBuiltAppProductionDependencies returns the AppProductionDependency objects for the application dependencies only
func GetBuiltAppProductionDependencies(appContext *context.AppContext) map[string]AppProductionDependency {
	result := map[string]AppProductionDependency{}
	for _, dependency := range appContext.Config.Production.Dependencies {
		builtDependency := NewAppProductionDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(appDir string, serviceRole string, appConfig types.AppConfig) error {
	if appConfig.Services == nil {
		appConfig.Services = map[string]types.ServiceSource{}
	}
	appConfig.Services[serviceRole] = types.ServiceSource{Location: fmt.Sprintf("./%s", serviceRole)}
	bytes, err := yaml.Marshal(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal application.yml")
	}
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), bytes, 0777)
}
