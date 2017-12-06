package types

// ServiceDevelopmentConfig represents development specific configuration for a service
type ServiceDevelopmentConfig struct {
	HealthCheck string            `yaml:"health-check,omitempty"`
	Port        string            `yaml:",omitempty"`
	Scripts     map[string]string `yaml:",omitempty"`
}
