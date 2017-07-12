package types

// DockerCompose represents the docker compose object
type DockerCompose struct {
	Version  string
	Services map[string]Service
}

// Service represents the service object in DockerCompose object
type Service struct {
	Image         string      `yaml:",omitempty"`
	Build         string      `yaml:",omitempty"`
	Command       string      `yaml:",omitempty"`
	ContainerName string      `yaml:",omitempty"`
	Ports         []string    `yaml:",omitempty"`
	Volumes       []string    `yaml:",omitempty"`
	Links         []string    `yaml:",omitempty"`
	Environment   Environment `yaml:",omitempty"`
	DependsOn     []string    `yaml:",omitempty"`
}

// Environment represents the environment variables of the service
type Environment struct {
	ServiceRoutes []ServiceRoute `yaml:"SERVICE_ROUTES,omitempty"`
}

// ServiceRoute represents the route that the service uses to communicate
// with exocom
type ServiceRoute struct {
	Role      string   `yaml:",omitempty"`
	Receives  []string `yaml:",omitempty"`
	Sends     []string `yaml:",omitempty"`
	Namespace string   `yaml:",omitempty"`
}
