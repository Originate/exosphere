package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GetBuiltLocalAppDependencies returns the LocalAppDependency objects for application dependencies only
func GetBuiltLocalAppDependencies(appContext *context.AppContext) map[string]*LocalAppDependency {
	result := map[string]*LocalAppDependency{}
	for name, dependency := range appContext.Config.Local.Dependencies {
		builtDependency := NewLocalAppDependency(name, dependency, appContext)
		result[name] = builtDependency
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

// GetBuiltLocalServiceDependencies returns the dependencies for a single service
func GetBuiltLocalServiceDependencies(serviceConfig types.ServiceConfig, appContext *context.AppContext) map[string]*LocalAppDependency {
	result := map[string]*LocalAppDependency{}
	for name, dependency := range serviceConfig.Local.Dependencies {
		builtDependency := NewLocalAppDependency(name, dependency, appContext)
		result[name] = builtDependency
	}
	return result
}

// GetAllRemoteDependencies returns all remote dependencies
func GetAllRemoteDependencies(appContext *context.AppContext) map[string]types.RemoteDependency {
	result := map[string]types.RemoteDependency{}
	for dependencyName, dependency := range appContext.Config.Remote.Dependencies {
		result[dependencyName] = dependency
	}
	for _, serviceContext := range appContext.ServiceContexts {
		for dependencyName, dependency := range serviceContext.Config.Remote.Dependencies {
			result[dependencyName] = dependency
		}
	}
	return result
}

// GetSortedRemoteDependencyNames returns all remote dependency names in alphabetical order
func GetSortedRemoteDependencyNames(appContext *context.AppContext) []string {
	result := []string{}
	for k := range appContext.Config.Remote.Dependencies {
		result = append(result, k)
	}
	for _, serviceContext := range appContext.ServiceContexts {
		for k := range serviceContext.Config.Remote.Dependencies {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}
