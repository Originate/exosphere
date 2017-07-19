package appDependencyHelpers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

type genericDependency struct {
	config    types.Dependency
	appConfig types.AppConfig
}

// GetContainerName returns the container name for the dependency
func (dependency genericDependency) GetContainerName() string {
	return dependency.config.Name + dependency.config.Version
}

// GetEnvVariables returns the environment variables for the depedency
func (dependency genericDependency) GetEnvVariables() map[string]string {
	return dependency.config.Config.DependencyEnvironment
}

// GetOnlineText returns the online text for the dependency
func (dependency genericDependency) GetOnlineText() string {
	return dependency.config.Config.OnlineText
}

// GetDockerConfig returns docker configuration for the dependency
func (dependency genericDependency) GetDockerConfig() (types.DockerConfig, error) {
	renderedVolumes, err := dependency.getRenderedVolumes()
	if err != nil {
		return types.DockerConfig{}, err
	}
	return types.DockerConfig{
		Image:         fmt.Sprintf("%s:%s", dependency.config.Name, dependency.config.Version),
		ContainerName: dependency.GetContainerName(),
		Ports:         dependency.config.Config.Ports,
		Volumes:       renderedVolumes,
	}, nil
}

func (dependency genericDependency) getRenderedVolumes() ([]string, error) {
	homeDir, err := osHelpers.GetUserHomeDir()
	if err != nil {
		return []string{}, err
	}
	dataPath := path.Join(homeDir, ".exosphere", dependency.appConfig.Name, dependency.config.Name, "data")
	renderedVolumes := []string{}
	if err := os.MkdirAll(dataPath, 0777); err != nil { //nolint gas
		return renderedVolumes, errors.Wrap(err, "Failed to create the necessary directories for the volumes")
	}
	for _, volume := range dependency.config.Config.Volumes {
		renderedVolumes = append(renderedVolumes, strings.Replace(volume, "{{EXO_DATA_PATH}}", dataPath, -1))
	}
	return renderedVolumes, nil
}
