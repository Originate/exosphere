package types

// ServiceDevelopmentConfig represents development specific configuration for a service
type ServiceDevelopmentConfig struct {
	Scripts map[string]string `yaml:",omitempty"`
	Port    string            `yaml:",omitempty"`
}
