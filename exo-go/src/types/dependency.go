package types

import (
	"fmt"
	"os"
)

// Dependency represents a dependency of the application
type Dependency struct {
	Name    string
	Version string
	Silent  bool             `yaml:",omitempty"`
	Config  DependencyConfig `yaml:",omitempty"`
}

// DependencyConfig represents the configuration of an application
type DependencyConfig struct {
	Ports                 []string               `yaml:",omitempty"`
	Volumes               []string               `yaml:",omitempty"`
	OnlineText            string                 `yaml:"online-text,omitempty"`
	DependencyEnvironment map[string]interface{} `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]interface{} `yaml:"service-environment,omitempty"`
}

// GetEnvVariables returns a map that maps the environment variables of
// the depedency to their values
func (dependency *Dependency) GetEnvVariables() map[string]interface{} {
	switch dependency.Name {
	case "exocom":
		port := os.Getenv("EXOCOM_PORT")
		if len(port) == 0 {
			port = "80"
		}
		return map[string]interface{}{"EXOCOM_PORT": port}
	case "nats":
		return map[string]interface{}{"NATS_HOST": dependency.GetContainerName()}
	default:
		return dependency.Config.DependencyEnvironment
	}
}

// GetContainerName returns the container name for the dependency
func (dependency *Dependency) GetContainerName() string {
	return fmt.Sprintf("%s%s", dependency.Name, dependency.Version)
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
