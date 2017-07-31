package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/pkg/errors"
	"github.com/segmentio/go-prompt"
	"gopkg.in/yaml.v2"
)

// GetAllDependencyNames returns the container names (name+version) of all application
// and service dependencies
func GetAllDependencyNames(appDir string, appConfig types.AppConfig) ([]string, error) {
	result := []string{}
	for _, dependency := range appConfig.Dependencies {
		containerNames := dependency.Name + dependency.Version
		if !util.DoesStringArrayContain(result, containerNames) {
			result = append(result, containerNames)
		}
	}
	serviceConfigs, err := GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return result, err
	}
	for _, serviceConfig := range serviceConfigs {
		for _, containerName := range GetServiceDependencies(serviceConfig, appConfig) {
			if !util.DoesStringArrayContain(result, containerName) {
				result = append(result, containerName)
			}
		}
	}
	return result, nil
}

// GetEnvironmentVariables returns the environment variables of
// all dependencies listed in appConfig
func GetEnvironmentVariables(appConfig types.AppConfig, appDir, homeDir string) map[string]string {
	result := map[string]string{}
	for _, dependency := range appConfig.Dependencies {
		for variable, value := range NewAppDependency(dependency, appConfig, appDir, homeDir).GetEnvVariables() {
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
