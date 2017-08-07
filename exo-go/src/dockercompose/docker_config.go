package dockercompose

// DockerConfig represents the configuration of a service/dependency as provided in
// docker-compose.yml
type DockerConfig struct {
	Image         string            `yaml:",omitempty"`
	Build         map[string]string `yaml:",omitempty"`
	Command       string            `yaml:",omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
	Ports         []string          `yaml:",omitempty"`
	Volumes       []string          `yaml:",omitempty"`
	Links         []string          `yaml:",omitempty"`
	Environment   map[string]string `yaml:",omitempty"`
	DependsOn     []string          `yaml:"depends_on,omitempty"`
}
