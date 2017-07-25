package appConfigHelpers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
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
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return result, err
	}
	for _, serviceConfig := range serviceConfigs {
		for _, containerName := range serviceConfigHelpers.GetServiceDependencies(serviceConfig, appConfig) {
			if !util.DoesStringArrayContain(result, containerName) {
				result = append(result, containerName)
			}
		}
	}
	return result, nil
}

// GetAppConfig reads application.yml and returns the appConfig object
func GetAppConfig(appDir string) (result types.AppConfig, err error) {
	yamlFile, err := ioutil.ReadFile(path.Join(appDir, "application.yml"))
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "Failed to unmarshal application.yml")
	}
	return result, nil
}

// GetDependencyNames returns the names of all dependencies listed in appConfig
func GetDependencyNames(appConfig types.AppConfig) []string {
	result := []string{}
	for _, dependency := range appConfig.Dependencies {
		result = append(result, dependency.Name)
	}
	return result
}

// GetEnvironmentVariables returns the environment variables of
// all dependencies listed in appConfig
func GetEnvironmentVariables(appConfig types.AppConfig, appDir, homeDir string) map[string]string {
	result := map[string]string{}
	for _, dependency := range appConfig.Dependencies {
		for variable, value := range appDependencyHelpers.Build(dependency, appConfig, appDir, homeDir).GetEnvVariables() {
			result[variable] = value
		}
	}
	return result
}

// GetServiceNames returns the service names for the given services
func GetServiceNames(services types.Services) []string {
	result := []string{}
	for serviceName := range services.Private {
		result = append(result, serviceName)
	}
	for serviceName := range services.Public {
		result = append(result, serviceName)
	}
	return result
}

// GetSilencedDependencyNames returns the names of dependencies that are
// configured as silent
func GetSilencedDependencyNames(appConfig types.AppConfig) []string {
	result := []string{}
	for _, dependency := range appConfig.Dependencies {
		if dependency.Silent {
			result = append(result, dependency.Name)
		}
	}
	return result
}

// GetSilencedServiceNames returns the names of services that are configured
// as silent
func GetSilencedServiceNames(services types.Services) []string {
	result := []string{}
	for serviceName, serviceConfig := range services.Private {
		if serviceConfig.Silent {
			result = append(result, serviceName)
		}
	}
	for serviceName, serviceConfig := range services.Public {
		if serviceConfig.Silent {
			result = append(result, serviceName)
		}
	}
	return result
}

// GetServiceProtectionLevels returns a map containing service names to their protection level
func GetServiceProtectionLevels(appConfig types.AppConfig) map[string]string {
	result := make(map[string]string)
	for serviceName := range appConfig.Services.Private {
		result[serviceName] = "private"
	}
	for serviceName := range appConfig.Services.Public {
		result[serviceName] = "public"
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

// VerifyServiceDoesNotExist returns an error if the service serviceRole already
// exists in existingServices, and return nil otherwise.
func VerifyServiceDoesNotExist(serviceRole string, existingServices []string) error {
	if util.DoesStringArrayContain(existingServices, serviceRole) {
		return fmt.Errorf(`Service %v already exists in this application`, serviceRole)
	}
	return nil
}
