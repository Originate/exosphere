package types

type EnvVars struct {
	Default     map[string]string `yaml:",omitempty"`
	Development map[string]string `yaml:",omitempty"`
	Production  map[string]string `yaml:",omitempty"`
	Secrets     []string          `yaml:",omitempty"`
}
