package types

import (
	"fmt"
	"reflect"
)

// ServiceProductionConfig represents production specific configuration for an application
type ServiceProductionConfig struct {
	Port        string `yaml:"port,omitempty"`
	HealthCheck string `yaml:"health-check,omitempty"`
}

// ValidateProductionFields validates production fields
func (p ServiceProductionConfig) ValidateProductionFields(serviceLocation, protectionLevel string) error {
	if protectionLevel == ServiceTypePublic {
		requiredFields := []string{"Port", "HealthCheck"}
		for _, field := range requiredFields {
			value := reflect.ValueOf(p).FieldByName(field).String()
			if value == "" {
				return fmt.Errorf("%s/service.yml missing required field 'production.%s'", serviceLocation, field)
			}
		}
	}
	return nil
}
