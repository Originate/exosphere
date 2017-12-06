package types

import (
	"fmt"
	"reflect"
)

// ServiceRemoteConfig represents production specific configuration for an application
type ServiceRemoteConfig struct {
	Dependencies []RemoteDependency
	Environment  map[string]string
	Secrets      []string
	URL          string `yaml:"url,omitempty"`
	CPU          string `yaml:"cpu,omitempty"`
	Memory       string `yaml:"memory,omitempty"`
	HealthCheck  string `yaml:"health-check,omitempty"`
}

// ValidateRemoteFields validates that service.yml contiains the required fields
func (r ServiceRemoteConfig) ValidateRemoteFields(serviceLocation, protectionLevel string) error {
	requiredPublicFields := []string{"URL", "CPU", "Memory", "HealthCheck"}
	requiredWorkerFields := []string{"CPU", "Memory"}

	requiredFields := []string{}
	switch protectionLevel {
	case ServiceTypePublic:
		requiredFields = requiredPublicFields
	case ServiceTypeWorker:
		requiredFields = requiredWorkerFields
	}
	for _, field := range requiredFields {
		value := reflect.ValueOf(r).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("%s/service.yml missing required field 'remote.%s'", serviceLocation, field)
		}
	}
	return nil
}
