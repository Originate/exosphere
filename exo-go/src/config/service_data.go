package config

// ServiceData represents the service info as provided in application.yml
type ServiceData struct {
	Location    string `yaml:",omitempty"`
	DockerImage string `yaml:"docker-image,omitempty"`
	NameSpace   string `yaml:",omitempty"`
	Silent      bool   `yaml:",omitempty"`
}
