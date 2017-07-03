package types

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
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
	Ports                 []string          `yaml:",omitempty"`
	Volumes               []string          `yaml:",omitempty"`
	OnlineText            string            `yaml:"online-text,omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
}

func (dependency *Dependency) GetContainerName() string {
	return fmt.Sprintf("%s%s", dependency.Name, dependency.Verison)
}

func (dependency *Dependency) GetDockerConfig(appConfig AppConfig, done interface{}) {
	config := map[string]string{"image": fmt.Sprintf("%s:%s", dependency.Name, dependency.Verison), "container_name": dependency.GetContainerName(), "ports": dependency.Config.Ports, "volumes": "TODO"}
}

func (dependency *Dependency) renderVolumes(volumes []string, appName string) []string {
	dataPath := path.Join(user.Current(), ".exosphere", appName, dependency.Name, "data")
	if err := os.MkdirAll(dataPath, 0777); err != nil {
		log.Fatalf("Failed to create %s")
	}

}
