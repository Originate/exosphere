package types

// ServiceConfig represents the configuration of a service
type ServiceConfig struct {
	Location    string `yaml:",omitempty"`
	DockerImage string `yaml:"docker-image,omitempty"`
	NameSpace   string `yaml:",omitempty"`
	Silent      bool   `yaml:",omitempty"`
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
