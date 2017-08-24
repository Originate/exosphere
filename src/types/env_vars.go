package types

// EnvVars is the environment key listed in a ServiceConfig
type EnvVars struct {
	Default     map[string]string `yaml:",omitempty"`
	Development map[string]string `yaml:",omitempty"`
	Production  map[string]string `yaml:",omitempty"`
	Secrets     []string          `yaml:",omitempty"`
}
