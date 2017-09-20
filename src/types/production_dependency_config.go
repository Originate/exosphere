package types

// ProductionDependencyConfig represents a production dependency
type ProductionDependencyConfig struct {
	Config  ProductionDependencyConfigOptions `yaml:",omitempty"`
	Silent  bool                              `yaml:",omitempty"`
	Name    string
	Version string
}
