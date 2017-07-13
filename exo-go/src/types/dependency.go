package types

import "os"

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

// GetContainerName returns the container name for the dependency
func (dependency *Dependency) GetContainerName() string {
	return dependency.Name + dependency.Version
}

// GetEnvVariables returns a map that maps the environment variables of
// the depedency to their values
func (dependency *Dependency) GetEnvVariables() map[string]string {
	switch dependency.Name {
	case "exocom":
		port := os.Getenv("EXOCOM_PORT")
		if len(port) == 0 {
			port = "80"
		}
		return map[string]string{"EXOCOM_PORT": port}
	case "nats":
		return map[string]string{"NATS_HOST": dependency.GetContainerName()}
	default:
		return dependency.Config.DependencyEnvironment
	}
}

// GetOnlineText returns the online text for the dependency
func (dependency *Dependency) GetOnlineText() string {
	switch dependency.Name {
	case "exocom":
		return "ExoCom WebSocket listener online"
	case "nats":
		return "Listening for route connections"
	default:
		return dependency.Config.OnlineText
	}
}
