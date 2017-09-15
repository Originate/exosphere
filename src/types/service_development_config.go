package types

// ServiceDevelopmentConfig represents development specific configuration for a service
type ServiceDevelopmentConfig struct {
	Dependencies []DependencyConfig
	Scripts      map[string]string `yaml:",omitempty"`
}
