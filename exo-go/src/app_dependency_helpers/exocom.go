package appDependencyHelpers

import (
	"fmt"
	"os"

	"github.com/Originate/exosphere/exo-go/src/types"
)

type exocomDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

// GetContainerName returns the container name for the dependency
func (exocom exocomDependency) GetContainerName() string {
	return exocom.config.Name + exocom.config.Version
}

// GetEnvVariables returns the environment variables for the depedency
func (exocom exocomDependency) GetEnvVariables() map[string]string {
	port := os.Getenv("EXOCOM_PORT")
	if len(port) == 0 {
		port = "80"
	}
	return map[string]string{"EXOCOM_PORT": port}
}

// GetOnlineText returns the online text for the exocom
func (exocom exocomDependency) GetOnlineText() string {
	return "ExoCom WebSocket listener online"
}

// GetDockerConfig returns docker configuration for the dependency
func (exocom exocomDependency) GetDockerConfig() (types.DockerConfig, error) {
	serviceRoutes := "{}" // TODO: get json
	return types.DockerConfig{
		ContainerName: exocom.GetContainerName(),
		Image:         fmt.Sprintf("originate/exocom:%s", exocom.config.Version),
		Command:       "bin/exocom",
		Environment: map[string]string{
			"ROLE":           "exocom",
			"PORT":           "$EXOCOM_PORT",
			"SERVICE_ROUTES": serviceRoutes,
		},
	}, nil
}
