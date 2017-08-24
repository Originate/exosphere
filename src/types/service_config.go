package types

import "github.com/Originate/exosphere/src/util"

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type            string                 `yaml:",omitempty"`
	Description     string                 `yaml:",omitempty"`
	Author          string                 `yaml:",omitempty"`
	Setup           string                 `yaml:",omitempty"`
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
