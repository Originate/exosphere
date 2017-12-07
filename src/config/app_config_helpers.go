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

// GetBuiltLocalAppDependencies returns the LocalAppDependency objects for application dependencies only
func GetBuiltLocalAppDependencies(appContext *context.AppContext) map[string]LocalAppDependency {
	result := map[string]LocalAppDependency{}
	for name, dependency := range appContext.Config.Local.Dependencies {
		builtDependency := NewLocalAppDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}

// GetBuiltRemoteDependencies returns the RemoteAppDependency objects for the application and service
// prod dependencies of the entire application
func GetBuiltRemoteDependencies(appContext *context.AppContext) map[string]RemoteAppDependency {
	result := GetBuiltRemoteAppDependencies(appContext)
	for _, serviceContext := range appContext.ServiceContexts {
		for dependencyName, builtDependency := range GetBuiltRemoteServiceDependencies(serviceContext.Config, appContext) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetBuiltRemoteAppDependencies returns the RemoteAppDependency objects for the application dependencies only
func GetBuiltRemoteAppDependencies(appContext *context.AppContext) map[string]RemoteAppDependency {
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
