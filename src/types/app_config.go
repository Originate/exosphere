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
	Local       LocalConfig     `yaml:",omitempty"`
	Remote      AppRemoteConfig `yaml:",omitempty"`
	Services    map[string]ServiceSource
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
	for _, serviceSource := range result.Services {
		serviceSource.DependencyData.StringifyMapKeys()
	}
	return result, result.validateAppConfig()
}

// GetSortedServiceRoles returns the service roles listed in application.yml sorted alphabetically
func (a AppConfig) GetSortedServiceRoles() []string {
	result := []string{}
	for serviceRole := range a.Services {
		result = append(result, serviceRole)
	}
	sort.Strings(result)
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

func (a AppConfig) validateAppConfig() error {
	appNameRegex := regexp.MustCompile("^[a-z0-9]+(-[a-z0-9]+)*$")
	if !appNameRegex.MatchString(a.Name) {
		return fmt.Errorf("The 'name' field '%s' in application.yml is invalid. Only lowercase alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-z0-9]+(-[a-z0-9]+)*$/", a.Name)
	}
	var err error
	serviceRoleRegex := regexp.MustCompile("^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$")
	for serviceRole := range a.Services {
		if !serviceRoleRegex.MatchString(serviceRole) {
			err = fmt.Errorf("The service key 'services.%s' in application.yml is invalid. Only alphanumeric character(s) separated by a single hyphen are allowed. Must match regex: /^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$/", serviceRole)
		}
	}
	return err
}
