package types

// EnvVars is the environment key listed in a ServiceConfig
type EnvVars struct {
	Default map[string]string `yaml:",omitempty"`
	Local   map[string]string `yaml:",omitempty"`
	Remote  map[string]string `yaml:",omitempty"`
	Secrets []string          `yaml:",omitempty"`
}
