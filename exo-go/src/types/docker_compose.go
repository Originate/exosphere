package types

// DockerCompose represents the docker compose object
type DockerCompose struct {
	Version  string
	Services map[string]Service
}

// Service represents the configuration of a service as provided in
// docker-compose.yml
type Service struct {
	Image         string            `yaml:",omitempty"`
	Build         string            `yaml:",omitempty"`
	Command       string            `yaml:",omitempty"`
	ContainerName string            `yaml:",omitempty"`
	Ports         []string          `yaml:",omitempty"`
	Volumes       []string          `yaml:",omitempty"`
	Links         []string          `yaml:",omitempty"`
	Environment   map[string]string `yaml:",omitempty"`
	DependsOn     []string          `yaml:",omitempty"`
}
