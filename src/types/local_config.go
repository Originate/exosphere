package types

// LocalConfig represents development specific configuration for an application
type LocalConfig struct {
	Dependencies         map[string]LocalDependency
	EnvironmentVariables map[string]string `yaml:"environment-variables,omitempty"`
	Secrets              []string          `yaml:",omitempty"`
}
