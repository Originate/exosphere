package types

// DevelopmentDependencyConfig represents a development dependency
type DevelopmentDependencyConfig struct {
	Config  DevelopmentDependencyConfigOptions `yaml:",omitempty"`
	Name    string
	Version string
}
