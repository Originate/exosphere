package application

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// Initializer sets up the app
type Initializer struct {
	AppConfig                types.AppConfig
	DockerComposeConfig      types.DockerCompose
	DockerComposeProjectName string
	AppDir                   string
	HomeDir                  string
	BuildMode                composebuilder.BuildMode
	logger                   *util.Logger
}

// NewInitializer is Initializer's constructor
func NewInitializer(appConfig types.AppConfig, logger *util.Logger, appDir, homeDir, dockerComposeProjectName string, mode composebuilder.BuildMode) (*Initializer, error) {
	appSetup := &Initializer{
		AppConfig:                appConfig,
		DockerComposeConfig:      types.DockerCompose{Version: "3"},
		DockerComposeProjectName: dockerComposeProjectName,
		AppDir:    appDir,
		HomeDir:   homeDir,
		BuildMode: mode,
		logger:    logger,
	}
	return appSetup, nil
}

func (i *Initializer) renderDockerCompose(dockerComposeDir string) error {
	bytes, err := yaml.Marshal(i.DockerComposeConfig)
	if err != nil {
		return err
	}
	if err := util.CreateEmptyDirectory(dockerComposeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dockerComposeDir, "docker-compose.yml"), bytes, 0777)
}

func (i *Initializer) setupDockerImages(dockerComposeDir string) error {
	opts := compose.BaseOptions{
		DockerComposeDir: dockerComposeDir,
		Logger:           i.logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", i.DockerComposeProjectName)},
	}
	err := compose.PullAllImages(opts)
	if err != nil {
		return err
	}
	return compose.BuildAllImages(opts)
}

// BuildDockerCompose builds the docker compose file, returning the docker compose dir
func (i *Initializer) BuildDockerCompose() (string, error) {
	dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: i.AppConfig,
		AppDir:    i.AppDir,
		BuildMode: i.BuildMode,
		HomeDir:   i.HomeDir,
	})
	if err != nil {
		return "", err
	}
	i.DockerComposeConfig.Services = dockerConfigs
	dockerComposeDir := path.Join(i.AppDir, "tmp")
	if err := i.renderDockerCompose(dockerComposeDir); err != nil {
		return "", err
	}
	return dockerComposeDir, nil
}

// Initialize sets up the entire app and returns an error if any
func (i *Initializer) Initialize() error {
	dockerComposeDir, err := i.BuildDockerCompose()
	if err != nil {
		return err
	}
	return i.setupDockerImages(dockerComposeDir)
}
