package types

// ServiceRemoteEnvironment represents configuration for a particular environment
type ServiceRemoteEnvironment struct {
	EnvironmentVariables map[string]string `yaml:"environment-variables,omitempty"`
	Secrets              []string          `yaml:",omitempty"`
	URL                  string            `yaml:"url,omitempty"`
}
