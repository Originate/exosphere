package types

// DependencyConfig represents a dependency of an application
type DependencyConfig struct {
	Config  DependencyConfigOptions `yaml:",omitempty"`
	Silent  bool                    `yaml:",omitempty"`
	Name    string
	Version string
}
