package types

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

func (dependency *Dependency) GetEnvVariables() map[string]interface{} {
	switch dependency.Name {
	case "exocom":
		var port int
		var err error
		if len(os.Getenv("EXOCOM_PORT")) > 0 {
			port, err = strconv.Atoi(os.Getenv("EXOCOM_PORT"))
			if err != nil {
				log.Fatalf("EXOCOM _PORT must be an integer: %s", err)
			}
		} else {
			port = 80
		}
		return map[string]interface{}{"EXOCOM_PORT": port}
	case "nats":
		return map[string]interface{}{"NATS_HOST": dependency.GetContainerName()}
	default:
		return dependency.Config.DependencyEnvironment
	}
}

func (dependency *Dependency) GetContainerName() string {
	return fmt.Sprintf("%s%s", dependency.Name, dependency.Version)
}

// func (dependency *Dependency) GetContainerName() string {
// 	return fmt.Sprintf("%s%s", dependency.Name, dependency.Version)
// }

// func (dependency *Dependency) GetDockerConfig(appConfig AppConfig, done interface{}) {
// 	config := map[string]interface{}{"image": fmt.Sprintf("%s:%s", dependency.Name, dependency.Version), "container_name": dependency.GetContainerName(), "ports": dependency.Config.Ports, "volumes": "TODO"}
// }

// func (dependency *Dependency) renderVolumes(volumes []string, appName string) []string {
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatalf("Failed to get user home directory: %s", err)
// 	}
// 	dataPath := path.Join(usr.HomeDir, ".exosphere", appName, dependency.Name, "data")
// 	if err := os.MkdirAll(dataPath, 0777); err != nil {
// 		log.Fatalf("Failed to create %s")
// 	}
// }
