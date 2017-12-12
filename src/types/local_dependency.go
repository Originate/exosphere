package types

// LocalDependency represents a development dependency
type LocalDependency struct {
	Config LocalDependencyConfig `yaml:",omitempty"`
	Image  string
}
