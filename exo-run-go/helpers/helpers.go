package helpers

import (
	"io/ioutil"
	"log"

	"github.com/Originate/exosphere/exo-run-go/types"
	"gopkg.in/yaml.v2"
)

func GetAppConfig() types.AppConfig {
	yamlFile, err := ioutil.ReadFile("application.yml")
	var appConfig types.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		log.Fatalf("Failed to read application.yml: %s", err)
	}
	return appConfig
}

func GetExistingServices(appConfig types.AppConfig) ([]string, []string) {
	services := []string{}
	silencedServices := []string{}
	for service, serviceConfig := range appConfig.Services.Private {
		services = append(services, service)
		if serviceConfig.Silent {
			silencedServices = append(silencedServices, service)
		}
	}
	for service, serviceConfig := range appConfig.Services.Public {
		services = append(services, service)
		if serviceConfig.Silent {
			silencedServices = append(silencedServices, service)
		}
	}
	return services, silencedServices
}

func GetSilencedDependencies(appConfig types.AppConfig) []string {
	silencedDependencies := []string{}
	for _, dependency := range appConfig.Dependencies {
		if dependency.Silent {
			silencedDependencies = append(silencedDependencies, dependency.Name)
		}
	}
	return silencedDependencies
}

func Contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}
