package types

import (
	"fmt"

	"github.com/Originate/exosphere/src/util"
)

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type            string                 `yaml:",omitempty"`
	Description     string                 `yaml:",omitempty"`
	Author          string                 `yaml:",omitempty"`
	Startup         map[string]string      `yaml:",omitempty"`
	Restart         map[string]interface{} `yaml:",omitempty"`
	Tests           string                 `yaml:",omitempty"`
	ServiceMessages `yaml:"messages,omitempty"`
	Docker          DockerConfig       `yaml:",omitempty"`
	Dependencies    []DependencyConfig `yaml:",omitempty"`
	Environment     EnvVars            `yaml:",omitempty"`
	Production      map[string]string  `yaml:",omitempty"`
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

// ValidateProductionFields validates that service.yml contiains a production field
// and the required fields
func (s ServiceConfig) ValidateProductionFields(serviceData, protectionLevel string) error {
	requiredPublicFields := []string{"url", "cpu", "memory", "public-port", "health-check"}
	requiredPrivateFields := []string{"cpu", "memory", "public-port", "health-check"}
	requiredWorkerFields := []string{"cpu", "memory"}

	if s.Production == nil {
		return fmt.Errorf("%s/service.yml missing required section 'production'", serviceData)
	}

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
		if s.Production[field] == "" {
			return fmt.Errorf("%s/service.yml missing required field 'production.%s'", serviceData, field)
		}
	}
	return nil
}
