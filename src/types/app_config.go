package types

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"

	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name        string
	Description string
	Version     string
	Development AppDevelopmentConfig `yaml:",omitempty"`
	Production  AppProductionConfig  `yaml:",omitempty"`
	Services    Services
	Templates   map[string]string `yaml:",omitempty"`
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

// GetDevelopmentDependencyNames returns the names of all dev dependencies listed in appConfig
func (a AppConfig) GetDevelopmentDependencyNames() []string {
	result := []string{}
	for _, dependency := range a.Development.Dependencies {
		result = append(result, dependency.Name)
	}
	return result
}

// GetProductionDependencyNames returns the names of all prod dependencies listed in appConfig
func (a AppConfig) GetProductionDependencyNames() []string {
	result := []string{}
	for _, dependency := range a.Production.Dependencies {
		result = append(result, dependency.Name)
	}
	return result
}

// GetServiceData returns the configuration data listed under a service role in application.yml
func (a AppConfig) GetServiceData() map[string]ServiceData {
	result := make(map[string]ServiceData)
	a.forEachService(func(serviceType, serviceRole string, data ServiceData) {
		result[serviceRole] = data
	})
	return result
}

// GetSortedServiceRoles returns the service roles listed in application.yml sorted alphabetically
func (a AppConfig) GetSortedServiceRoles() []string {
	result := []string{}
	a.forEachService(func(serviceType, serviceRole string, data ServiceData) {
		result = append(result, serviceRole)
	})
	sort.Strings(result)
	return result
}

// GetServiceProtectionLevels returns a map containing service role to protection level
func (a AppConfig) GetServiceProtectionLevels() map[string]string {
	result := make(map[string]string)
	a.forEachService(func(serviceType, serviceRole string, data ServiceData) {
		result[serviceRole] = serviceType
	})
	return result
}

// GetSilencedDevelopmentDependencyNames returns the names of development dependencies that are
// configured as silent
func (a AppConfig) GetSilencedDevelopmentDependencyNames() []string {
	result := []string{}
	for _, dependency := range a.Development.Dependencies {
		if dependency.Silent {
			result = append(result, dependency.Name)
		}
	}
	return result
}

// GetSilencedServiceRoles returns the names of services that are configured
// as silent
func (a AppConfig) GetSilencedServiceRoles() []string {
	result := []string{}
	a.forEachService(func(serviceType, serviceRole string, data ServiceData) {
		if data.Silent {
			result = append(result, serviceRole)
		}
	})
	return result
}

// VerifyServiceRoleDoesNotExist returns an error if the serviceRole already
// exists in existingServices, and return nil otherwise.
func (a AppConfig) VerifyServiceRoleDoesNotExist(serviceRole string) error {
	if util.DoesStringArrayContain(a.GetSortedServiceRoles(), serviceRole) {
		return fmt.Errorf(`Service role '%v' already exists in this application`, serviceRole)
	}
	return nil
}

func (a AppConfig) forEachService(fn func(string, string, ServiceData)) {
	for serviceRole, data := range a.Services.Worker {
		fn("worker", serviceRole, data)
	}
	for serviceRole, data := range a.Services.Private {
		fn("private", serviceRole, data)
	}
	for serviceRole, data := range a.Services.Public {
		fn("public", serviceRole, data)
	}
}
