package types

// LocalConfig represents development specific configuration for an application
type LocalConfig struct {
	Dependencies map[string]LocalDependency
	Environment  map[string]string `yaml:",omitempty"`
	Secrets      []string          `yaml:",omitempty"`
}
