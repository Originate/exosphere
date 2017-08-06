package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/stringplus"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services
	Templates  map[string]string `yaml:",omitempty"`
	Production map[string]string `yaml:",omitempty"`
}

// NewAppConfig reads application.yml and returns the appConfig object
func NewAppConfig(appDir string) (result AppConfig, err error) {
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
func (a AppConfig) GetDependencyNames() []string {
	result := []string{}
	for _, dependency := range a.Dependencies {
		result = append(result, dependency.Name)
	}
	return result
}

// GetServiceData returns the configurations data for the given services
func (a AppConfig) GetServiceData() map[string]ServiceData {
	result := make(map[string]ServiceData)
	for serviceName, data := range a.Services.Private {
		result[serviceName] = data
	}
	for serviceName, data := range a.Services.Public {
		result[serviceName] = data
	}
	return result
}

// GetServiceNames returns the service names for the given services
func (a AppConfig) GetServiceNames() []string {
	result := []string{}
	for serviceName := range a.Services.Private {
		result = append(result, serviceName)
	}
	for serviceName := range a.Services.Public {
		result = append(result, serviceName)
	}
	return result
}

// GetServiceProtectionLevels returns a map containing service names to their protection level
func (a AppConfig) GetServiceProtectionLevels() map[string]string {
	result := make(map[string]string)
	for serviceName := range a.Services.Private {
		result[serviceName] = "private"
	}
	for serviceName := range a.Services.Public {
		result[serviceName] = "public"
	}
	return result
}

// GetSilencedDependencyNames returns the names of dependencies that are
// configured as silent
func (a AppConfig) GetSilencedDependencyNames() []string {
	result := []string{}
	for _, dependency := range a.Dependencies {
		if dependency.Silent {
			result = append(result, dependency.Name)
		}
	}
	return result
}

// GetSilencedServiceNames returns the names of services that are configured
// as silent
func (a AppConfig) GetSilencedServiceNames() []string {
	result := []string{}
	for serviceName, serviceConfig := range a.Services.Private {
		if serviceConfig.Silent {
			result = append(result, serviceName)
		}
	}
	for serviceName, serviceConfig := range a.Services.Public {
		if serviceConfig.Silent {
			result = append(result, serviceName)
		}
	}
	return result
}

// VerifyServiceDoesNotExist returns an error if the service serviceRole already
// exists in existingServices, and return nil otherwise.
func (a AppConfig) VerifyServiceDoesNotExist(serviceRole string) error {
	if stringplus.DoesStringArrayContain(a.GetServiceNames(), serviceRole) {
		return fmt.Errorf(`Service %v already exists in this application`, serviceRole)
	}
	return nil
}
