package types

// DependencyConfig represents a dependency of an application
type DependencyConfig struct {
	Config  DependencyConfigOptions
	Silent  bool
	Name    string
	Version string
}
