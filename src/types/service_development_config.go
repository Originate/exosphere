package types

// ServiceDevelopmentConfig represents development specific configuration for a service
type ServiceDevelopmentConfig struct {
	Dependencies []DevelopmentDependencyConfig
	Scripts      map[string]string `yaml:",omitempty"`
}
