package types

import (
	"fmt"
	"reflect"
)

// ServiceProductionConfig represents production specific configuration for an application
type ServiceProductionConfig struct {
	Dependencies []ProductionDependencyConfig
	URL          string `yaml:"url,omitempty"`
	CPU          string `yaml:"cpu,omitempty"`
	Memory       string `yaml:"memory,omitempty"`
	Port         string `yaml:"port,omitempty"`
	HealthCheck  string `yaml:"health-check,omitempty"`
}

// ValidateFields validates that service.yml contiains the required fields
func (p ServiceProductionConfig) ValidateFields(serviceLocation, protectionLevel string) error {
	requiredPublicFields := []string{"URL", "CPU", "Memory", "Port", "HealthCheck"}
	requiredPrivateFields := []string{"CPU", "Memory", "Port", "HealthCheck"}
	requiredWorkerFields := []string{"CPU", "Memory"}

	requiredFields := []string{}
	switch protectionLevel {
	case "public":
		requiredFields = requiredPublicFields
	case "private":
		requiredFields = requiredPrivateFields
	case "worker":
		requiredFields = requiredWorkerFields
	}
	for _, field := range requiredFields {
		value := reflect.ValueOf(p).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("%s/service.yml missing required field 'production.%s'", serviceLocation, field)
		}
	}
	return nil
}
