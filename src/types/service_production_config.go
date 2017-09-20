package types

import (
	"fmt"
	"reflect"
)

// ServiceProductionConfig represents production specific configuration for an application
type ServiceProductionConfig struct {
	Dependencies []ProductionDependencyConfig
	URL          string `yaml:",omitempty"`
	CPU          string `yaml:"url,omitempty"`
	Memory       string `yaml:"cpu,omitempty"`
	PublicPort   string `yaml:"public-port,omitempty"`
	HealthCheck  string `yaml:"health-check,omitempty"`
}

// ValidateFields validates that service.yml contiains the required fields
func (p ServiceProductionConfig) ValidateFields(serviceLocation, protectionLevel string) error {
	requiredPublicFields := []string{"URL", "CPU", "Memory", "PublicPort", "HealthCheck"}
	requiredPrivateFields := []string{"CPU", "Memory", "PublicPort", "HealthCheck"}
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
