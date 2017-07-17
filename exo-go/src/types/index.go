package types

// Dependency represents a dependency of the application
type Dependency struct {
	Name    string
	Version string
	Silent  bool             `yaml:",omitempty"`
	Config  DependencyConfig `yaml:",omitempty"`
}

// DependencyConfig represents the configuration of an application
type DependencyConfig struct {
	Ports                 []string          `yaml:",omitempty"`
	Volumes               []string          `yaml:",omitempty"`
	OnlineText            string            `yaml:"online-text,omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
}

// ServiceConfig represents the configuration of a service as provided in
// service.yml
type ServiceConfig struct {
	Type            string            `yaml:",omitempty"`
	Description     string            `yaml:",omitempty"`
	Author          string            `yaml:",omitempty"`
	Setup           string            `yaml:",omitempty"`
	Startup         map[string]string `yaml:",omitempty"`
	Production      map[string]string `yaml:",omitempty"`
	ServiceMessages `yaml:",omitempty"`
}

// ServiceData represents the service info as provided in application.yml
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

// Services represents the mapping of protection level to services
type Services struct {
	Public  map[string]ServiceData
	Private map[string]ServiceData
}

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services
	Templates map[string]string `yaml:",omitempty"`
}
