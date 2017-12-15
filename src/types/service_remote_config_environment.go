package types

// ServiceRemoteConfig represents production specific configuration for an application
type ServiceRemoteConfigEnvironment struct {
	Environment map[string]string `yaml:",omitempty"`
	Secrets     []string          `yaml:",omitempty"`
	URL         string            `yaml:"url,omitempty"`
}
