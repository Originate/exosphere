package types

// LocalDependencyConfig represents the configuration of a local dependency
type LocalDependencyConfig struct {
	Ports                 []string          `yaml:",omitempty"`
	Persist               []string          `yaml:",omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
}
