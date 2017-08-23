package types

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
	Environment     envVars            `yaml:",omitempty"`
	Production      map[string]string  `yaml:",omitempty"`
}

type envVars struct {
	Default     map[string]string `yaml:",omitempty"`
	Development map[string]string `yaml:",omitempty"`
	Production  map[string]string `yaml:",omitempty"`
	Secrets     []string          `yaml:",omitempty"`
}

// GetEnvVars compiles a service's environment variables
func (s ServiceConfig) GetEnvVars(environment string) (map[string]string, []string) {
	defaultVars := s.Environment.Default
	envVars := map[string]string{}
	switch environment {
	case "production":
		envVars = s.Environment.Production
	case "development":
		envVars = s.Environment.Development
	}
	for k, v := range envVars {
		defaultVars[k] = v
	}
	return defaultVars, s.Environment.Secrets
}
