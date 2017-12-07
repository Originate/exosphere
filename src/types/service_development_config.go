package types

// ServiceDevelopmentConfig represents development specific configuration for a service
type ServiceDevelopmentConfig struct {
	Port    string            `yaml:",omitempty"`
	Scripts map[string]string `yaml:",omitempty"`
}
