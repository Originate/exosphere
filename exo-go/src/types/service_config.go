package types

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type            string                 `yaml:",omitempty"`
	Description     string                 `yaml:",omitempty"`
	Author          string                 `yaml:",omitempty"`
	Setup           string                 `yaml:",omitempty"`
	Startup         map[string]string      `yaml:",omitempty"`
	Restart         map[string]interface{} `yaml:",omitempty"`
	ServiceMessages `yaml:"messages,omitempty"`
	Docker          map[string]interface{} `yaml:",omitempty"`
	Production      map[string]string      `yaml:",omitempty"`
}
