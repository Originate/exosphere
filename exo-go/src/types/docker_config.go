package types

// DockerConfig represents the configuration of a service/dependency as provided in
// docker-compose.yml
type DockerConfig struct {
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
