package types

// ServiceRemoteEnvironment represents configuration for a particular environment
type ServiceRemoteEnvironment struct {
	Environment map[string]string `yaml:",omitempty"`
	Secrets     []string          `yaml:",omitempty"`
	URL         string            `yaml:"url,omitempty"`
}
