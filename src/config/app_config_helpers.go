package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetBuiltLocalAppDependencies returns the LocalAppDependency objects for application dependencies only
func GetBuiltLocalAppDependencies(appContext *types.AppContext) map[string]LocalAppDependency {
	result := map[string]LocalAppDependency{}
	for _, dependency := range appContext.Config.Local.Dependencies {
		builtDependency := NewLocalAppDependency(dependency, appContext)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetBuiltRemoteDependencies returns the RemoteAppDependency objects for the application and service
// prod dependencies of the entire application
func GetBuiltRemoteDependencies(appContext *types.AppContext, serviceConfigs map[string]types.ServiceConfig) map[string]RemoteAppDependency {
	result := GetBuiltRemoteAppDependencies(appContext)
	for _, serviceConfig := range serviceConfigs {
		for dependencyName, builtDependency := range GetBuiltRemoteServiceDependencies(serviceConfig, appContext) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetBuiltRemoteAppDependencies returns the RemoteAppDependency objects for the application dependencies only
func GetBuiltRemoteAppDependencies(appContext *types.AppContext) map[string]RemoteAppDependency {
	result := map[string]RemoteAppDependency{}
	for _, dependency := range appContext.Config.Remote.Dependencies {
		builtDependency := NewRemoteAppDependency(dependency, appContext)
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
