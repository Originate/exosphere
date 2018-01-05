package types

// LocalDependency represents a development dependency
type LocalDependency struct {
	Image                string
	Persist              []string          `yaml:",omitempty"`
	EnvironmentVariables map[string]string `yaml:"environment-variables,omitempty"`
	Secrets              []string          `yaml:",omitempty"`
}
