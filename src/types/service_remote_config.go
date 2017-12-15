package types

import (
	"fmt"
	"reflect"
)

// ServiceRemoteConfig represents production specific configuration for an application
type ServiceRemoteConfig struct {
	Dependencies map[string]RemoteDependency
	Environments map[string]ServiceRemoteConfigEnvironment
	CPU          string `yaml:"cpu,omitempty"`
	Memory       string `yaml:"memory,omitempty"`
}

// ValidateRemoteFields validates that service.yml contiains the required fields
func (r ServiceRemoteConfig) ValidateRemoteFields(serviceLocation, protectionLevel, remoteEnvironmentID string) error {
	requiredFields := []string{"CPU", "Memory"}
	for _, field := range requiredFields {
		value := reflect.ValueOf(r).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("%s/service.yml missing required field 'remote.%s'", serviceLocation, field)
		}
	}
	if protectionLevel == ServiceTypePublic && r.Environments[remoteEnvironmentID].URL == "" {
		return fmt.Errorf("%s/service.yml missing required field 'remote.environments.%s.url'", serviceLocation, remoteEnvironmentID)
	}
	return nil
}
