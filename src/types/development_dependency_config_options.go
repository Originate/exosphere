package types

// DevelopmentDependencyConfigOptions represents the configuration of a dependency
type DevelopmentDependencyConfigOptions struct {
	Ports                 []string          `yaml:",omitempty"`
	Persist               []string          `yaml:",omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
}
