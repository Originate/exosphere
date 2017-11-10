package types

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type            string `yaml:",omitempty"`
	Description     string `yaml:",omitempty"`
	Author          string `yaml:",omitempty"`
	ServiceMessages `yaml:"messages,omitempty"`
	Docker          DockerConfig             `yaml:",omitempty"`
	Environment     EnvVars                  `yaml:",omitempty"`
	Development     ServiceDevelopmentConfig `yaml:",omitempty"`
	Production      ServiceProductionConfig  `yaml:",omitempty"`
}

// NewServiceConfig returns a validated ServiceConfig object given the app directory path
// and the directory name of a service
func NewServiceConfig(appDir, serviceDirName string) (ServiceConfig, error) {
	var serviceConfig ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(appDir, serviceDirName, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the internal service '%s'", serviceDirName))
	}
	return serviceConfig, serviceConfig.ValidateServiceConfig()
}

// GetEnvVars compiles a service's environment variables
// It overwrites default variables with environemnt specific ones,
// returning a map of public env vars and a list of private env var keys
func (s ServiceConfig) GetEnvVars(environment string) (map[string]string, []string) {
	result := map[string]string{}
	util.Merge(result, s.Environment.Default)
	envVars := map[string]string{}
	switch environment {
	case "production":
		envVars = s.Environment.Production
	case "development":
		envVars = s.Environment.Development
	}
	util.Merge(result, envVars)
	return result, s.Environment.Secrets
}

// ValidateServiceConfig validates a ServiceConfig object
func (s ServiceConfig) ValidateServiceConfig() error {
	validTypes := []string{"public", "worker"}
	if !util.DoesStringArrayContain(validTypes, s.Type) {
		fmt.Println("err~~~~~~~~~~~~")
		fmt.Println(s.Type)
		return fmt.Errorf("Invalid value '%s' in service.yml field 'type'. Must be one of: %s", s.Type, strings.Join(validTypes, ", "))
	}
	return nil
}
