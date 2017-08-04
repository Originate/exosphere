package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
	prompt "github.com/segmentio/go-prompt"
	"gopkg.in/yaml.v2"
)

// GetAllBuiltDependencies returns the AppDependency objects for the
// dependencies of the entire application
func GetAllBuiltDependencies(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir, homeDir string) map[string]AppDependency {
	result := GetAppBuiltDependencies(appConfig, appDir, homeDir)
	for _, serviceConfig := range serviceConfigs {
		for dependencyName, builtDependency := range GetServiceBuiltDependencies(serviceConfig, appConfig, appDir, homeDir) {
			result[dependencyName] = builtDependency
		}
	}
	return result
}

// GetAppBuiltDependencies returns the AppDependency objects for the
// dependencies defined in the given appConfig
func GetAppBuiltDependencies(appConfig types.AppConfig, appDir, homeDir string) map[string]AppDependency {
	result := map[string]AppDependency{}
	for _, dependency := range appConfig.Dependencies {
		builtDependency := NewAppDependency(dependency, appConfig, appDir, homeDir)
		result[dependency.Name] = builtDependency
	}
	return result
}

// GetEnvironmentVariables returns the environment variables of
// all dependencies listed in appConfig
func GetEnvironmentVariables(appBuiltDependencies map[string]AppDependency) map[string]string {
	result := map[string]string{}
	for _, builtDependency := range appBuiltDependencies {
		for variable, value := range builtDependency.GetEnvVariables() {
			result[variable] = value
		}
	}
	return result
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(appDir string, serviceRole string, appConfig types.AppConfig) error {
	protectionLevels := []string{"public", "private"}
	switch protectionLevels[prompt.Choose("Protection Level:", protectionLevels)] {
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
	}
	bytes, err := yaml.Marshal(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal application.yml")
	}
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), bytes, 0777)
}
