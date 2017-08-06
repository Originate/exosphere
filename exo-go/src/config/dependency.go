package config

// Dependency represents a dependency of an application
type Dependency struct {
	Name    string
	Version string
	Silent  bool             `yaml:",omitempty"`
	Config  DependencyConfig `yaml:",omitempty"`
}
