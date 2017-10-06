package types

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
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
	return result, result.validateAppConfig()
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

// GetServiceData returns the configurations data for the given services
func (a AppConfig) GetServiceData() map[string]ServiceData {
	result := make(map[string]ServiceData)
	a.forEachService(func(serviceType, serviceName string, data ServiceData) {
		result[serviceName] = data
	})
	return result
}

// GetSortedServiceNames returns the service names for the given services sorted alphabetically
func (a AppConfig) GetSortedServiceNames() []string {
	result := []string{}
	a.forEachService(func(serviceType, serviceName string, data ServiceData) {
		result = append(result, serviceName)
	})
	sort.Strings(result)
	return result
}

// GetServiceProtectionLevels returns a map containing service names to their protection level
func (a AppConfig) GetServiceProtectionLevels() map[string]string {
	result := make(map[string]string)
	a.forEachService(func(serviceType, serviceName string, data ServiceData) {
		result[serviceName] = serviceType
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

// GetSilencedServiceNames returns the names of services that are configured
// as silent
func (a AppConfig) GetSilencedServiceNames() []string {
	result := []string{}
	a.forEachService(func(serviceType, serviceName string, data ServiceData) {
		if data.Silent {
			result = append(result, serviceName)
		}
	})
	return result
}

// VerifyServiceDoesNotExist returns an error if the service serviceRole already
// exists in existingServices, and return nil otherwise.
func (a AppConfig) VerifyServiceDoesNotExist(serviceRole string) error {
	if util.DoesStringArrayContain(a.GetSortedServiceNames(), serviceRole) {
		return fmt.Errorf(`Service %v already exists in this application`, serviceRole)
	}
	return nil
}

func (a AppConfig) forEachService(fn func(string, string, ServiceData)) {
	for serviceName, data := range a.Services.Worker {
		fn("worker", serviceName, data)
	}
	for serviceName, data := range a.Services.Private {
		fn("private", serviceName, data)
	}
	for serviceName, data := range a.Services.Public {
		fn("public", serviceName, data)
	}
}

func (a AppConfig) validateAppConfig() error {
	appNameRegex := regexp.MustCompile("^[a-z0-9]+(-[a-z0-9]+)*$")
	if !appNameRegex.MatchString(a.Name) {
		return errors.New("only lowercase alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-z0-9]+(-[a-z0-9]+)*$/")
	}
	var err error
	serviceRoleRegex := regexp.MustCompile("^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$")
	a.forEachService(func(serviceType, serviceRole string, data ServiceData) {
		if !serviceRoleRegex.MatchString(serviceRole) {
			err = errors.New("only alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$/")
		}
	})
	return err
}
