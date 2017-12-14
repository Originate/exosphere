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

// ServiceTypePublic is the value for the type field of a public service
const ServiceTypePublic = "public"

// ServiceTypeWorker is the value for the type field of a worker service
const ServiceTypeWorker = "worker"

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type           string                         `yaml:",omitempty"`
	Description    string                         `yaml:",omitempty"`
	Author         string                         `yaml:",omitempty"`
	DependencyData ServiceDependencyData          `yaml:"dependency-data,omitempty"`
	Docker         DockerConfig                   `yaml:",omitempty"`
	Development    ServiceDevelopmentConfig       `yaml:",omitempty"`
	Local          LocalConfig                    `yaml:",omitempty"`
	Production     ServiceProductionConfig        `yaml:",omitempty"`
	Remote         map[string]ServiceRemoteConfig `yaml:",omitempty"`
}

// NewServiceConfig returns a validated ServiceConfig object given the app directory path
// and the directory name of a service
func NewServiceConfig(serviceLocation string) (ServiceConfig, error) {
	var serviceConfig ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(serviceLocation, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the internal service '%s'", path.Base(serviceLocation)))
	}
	serviceConfig.DependencyData.StringifyMapKeys()
	return serviceConfig, serviceConfig.ValidateServiceConfig()
}

// ValidateServiceConfig validates a ServiceConfig object
func (s ServiceConfig) ValidateServiceConfig() error {
	validTypes := []string{ServiceTypePublic, ServiceTypeWorker}
	if !util.DoesStringArrayContain(validTypes, s.Type) {
		return fmt.Errorf("Invalid value '%s' in service.yml field 'type'. Must be one of: %s", s.Type, strings.Join(validTypes, ", "))
	}
	return nil
}

// ValidateDeployFields validates a serviceConfig for deployment
func (s ServiceConfig) ValidateDeployFields(serviceLocation, remoteID string) error {
	err := s.Production.ValidateProductionFields(serviceLocation, s.Type)
	if err != nil {
		return err
	}
	return s.Remote[remoteID].ValidateRemoteFields(serviceLocation, s.Type)
}
