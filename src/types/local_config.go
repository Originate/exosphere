package types

// LocalConfig represents development specific configuration for an application
type LocalConfig struct {
	Dependencies []LocalDependency
	Environment  map[string]string
	Secrets      []string
}
