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

// GetEnvironmentVariables returns environment variables ("VAR=VALUE")
// of the dependencies listed in appConfig
func GetEnvironmentVariables(appConfig types.AppConfig) map[string]string {
	envVars := map[string]string{}
	for _, dependency := range appConfig.Dependencies {
		for variable, value := range dependency.GetEnvVariables() {
			envVars[variable] = value
		}
	}
	return envVars
}

// GetDependencyNames returns dependency names listed in appConfig
func GetDependencyNames(appConfig types.AppConfig) []string {
	dependencyNames := []string{}
	for _, dependency := range appConfig.Dependencies {
		dependencyNames = append(dependencyNames, dependency.Name)
	}
	return dependencyNames
}

// GetServiceNames returns service names for the given services
func GetServiceNames(services types.Services) []string {
	serviceNames := []string{}
	for serviceName := range services.Private {
		serviceNames = append(serviceNames, serviceName)
	}
	for serviceName := range services.Public {
		serviceNames = append(serviceNames, serviceName)
	}
	return serviceNames
}

// GetSilencedDependencyNames returns dependency names that are configured as silent.
func GetSilencedDependencyNames(appConfig types.AppConfig) []string {
	silencedDependencyNames := []string{}
	for _, dependency := range appConfig.Dependencies {
		if dependency.Silent {
			silencedDependencyNames = append(silencedDependencyNames, dependency.Name)
		}
	}
	return silencedDependencyNames
}

// GetSilencedServiceNames returns service names that are configured as silent.
func GetSilencedServiceNames(services types.Services) []string {
	silencedServiceNames := []string{}
	for serviceName, serviceConfig := range services.Private {
		if serviceConfig.Silent {
			silencedServiceNames = append(silencedServiceNames, serviceName)
		}
	}
	for serviceName, serviceConfig := range services.Public {
		if serviceConfig.Silent {
			silencedServiceNames = append(silencedServiceNames, serviceName)
		}
	}
	return silencedServiceNames
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
