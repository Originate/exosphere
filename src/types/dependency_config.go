package types

// DependencyConfig represents a dependency of an application
type DependencyConfig struct {
	Name    string
	Version string
	Silent  bool                    `yaml:",omitempty"`
	Config  DependencyConfigOptions `yaml:",omitempty"`
}
