package types

// ServiceConfig represents the configuration of a service
type ServiceConfig struct {
	Location    string            `yaml:",omitempty"`
	DockerImage string            `yaml:"docker-image,omitempty"`
	NameSpace   string            `yaml:",omitempty"`
	Silent      bool              `yaml:",omitempty"`
	Type        string            `yaml:",omitempty"`
	Description string            `yaml:",omitempty"`
	Author      string            `yaml:",omitempty"`
	Setup       string            `yaml:",omitempty"`
	Startup     map[string]string `yaml:",omitempty"`
	Messages    `yaml:",omitempty"`
}

// Messages represents the messages that the service sends and receives
type Messages struct {
	Receives []string
	Sends    []string
}

// Services represents the mapping of protection level to services
type Services struct {
	Public  map[string]ServiceConfig
	Private map[string]ServiceConfig
}

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services
}
