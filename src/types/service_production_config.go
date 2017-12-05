package types

import "fmt"

// ServiceProductionConfig represents production specific configuration for an application
type ServiceProductionConfig struct {
	Port string `yaml:"port,omitempty"`
}

// ValidateProductionFields validates production fields
func (p ServiceProductionConfig) ValidateProductionFields(serviceLocation, protectionLevel string) error {
	if protectionLevel == "public" && p.Port == "" {
		return fmt.Errorf("%s/service.yml missing required field 'production.Port'", serviceLocation)
	}
	return nil
}
