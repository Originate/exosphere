package types

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services
	Templates map[string]string `yaml:",omitempty"`
}

// Services represents the mapping of protection level to services
type Services struct {
	Public  map[string]ServiceData
	Private map[string]ServiceData
}

// ServiceConfig represents the configuration of a service provided in
// service.yml
type ServiceConfig struct {
	Type            string            `yaml:",omitempty"`
	Description     string            `yaml:",omitempty"`
	Author          string            `yaml:",omitempty"`
	Setup           string            `yaml:",omitempty"`
	Startup         map[string]string `yaml:",omitempty"`
	ServiceMessages `yaml:",omitempty"`
}

// ServiceData contains service info provided in application.yml
type ServiceData struct {
	Location    string `yaml:",omitempty"`
	DockerImage string `yaml:"docker-image,omitempty"`
	NameSpace   string `yaml:",omitempty"`
	Silent      bool   `yaml:",omitempty"`
}

// ServiceMessages represents the messages that the service sends and receives
type ServiceMessages struct {
	Receives []string
	Sends    []string
}
