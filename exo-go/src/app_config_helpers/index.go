package appConfigHelpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/user_input_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"gopkg.in/yaml.v2"
)

// GetAppConfig reads application.yml and returns the appConfig object
func GetAppConfig() types.AppConfig {
	yamlFile, err := ioutil.ReadFile("application.yml")
	if err != nil {
		log.Fatalf("Failed to read application.yml: %s", err)
	}
	var appConfig types.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal application.yml: %s", err)
	}
	return appConfig
}

// GetEnvironmentVariables returns the environment variables of
// all dependencies listed in appConfig
func GetEnvironmentVariables(appConfig types.AppConfig) map[string]string {
	result := map[string]string{}
	for _, dependency := range appConfig.Dependencies {
		for variable, value := range dependency.GetEnvVariables() {
			result[variable] = value
		}
	}
	return result
}

// GetDependencyNames returns the names of all dependencies listed in appConfig
func GetDependencyNames(appConfig types.AppConfig) []string {
	result := []string{}
	for _, dependency := range appConfig.Dependencies {
		result = append(result, dependency.Name)
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

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(reader *bufio.Reader, serviceRole string, appConfig types.AppConfig) {
	switch userInputHelpers.Choose(reader, "Protection Level:", []string{"public", "private"}) {
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
		log.Fatalf("Failed to marshal application.yml: %s", err)
	}
	err = ioutil.WriteFile(path.Join("application.yml"), bytes, 0777)
	if err != nil {
		log.Fatalf("Failed to write application.yml: %s", err)
	}
}

// VerifyServiceDoesNotExist returns an error if the service serviceRole already
// exists in existingServices, and return nil otherwise.
func VerifyServiceDoesNotExist(serviceRole string, existingServices []string) error {
	if util.DoesStringArrayContain(existingServices, serviceRole) {
		return fmt.Errorf(`Service %v already exists in this application`, serviceRole)
	}
	return nil
}
