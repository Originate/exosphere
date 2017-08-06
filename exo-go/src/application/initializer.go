package application

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/dockercompose"
	"github.com/Originate/exosphere/exo-go/src/dockercomposebuilder"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/osplus"
)

// Initializer sets up the app
type Initializer struct {
	AppConfig            config.AppConfig
	BuiltAppDependencies map[string]config.AppDependency
	DockerComposeConfig  dockercompose.DockerCompose
	ServiceData          map[string]config.ServiceData
	ServiceConfigs       map[string]config.ServiceConfig
	AppDir               string
	HomeDir              string
	Logger               *logger.Logger
	logChannel           chan string
}

// NewInitializer is Initializer's constructor
func NewInitializer(appConfig config.AppConfig, logger *logger.Logger, appDir, homeDir string) (*Initializer, error) {
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Initializer{}, err
	}
	appSetup := &Initializer{
		AppConfig:            appConfig,
		BuiltAppDependencies: config.GetAppBuiltDependencies(appConfig, appDir, homeDir),
		DockerComposeConfig:  dockercompose.DockerCompose{Version: "3"},
		ServiceData:          appConfig.GetServiceData(),
		ServiceConfigs:       serviceConfigs,
		AppDir:               appDir,
		HomeDir:              homeDir,
		Logger:               logger,
		logChannel:           logger.GetLogChannel("exo-run"),
	}
	return appSetup, nil
}

func (i *Initializer) getAppDependenciesDockerConfigs() (dockercompose.DockerConfigs, error) {
	result := dockercompose.DockerConfigs{}
	for _, builtDependency := range i.BuiltAppDependencies {
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

// GetDockerConfigs returns the docker configs of all services and dependencies in
// the application
func (i *Initializer) GetDockerConfigs() (dockercompose.DockerConfigs, error) {
	dependencyDockerConfigs, err := i.getAppDependenciesDockerConfigs()
	if err != nil {
		return nil, err
	}
	serviceDockerConfigs, err := i.getServiceDockerConfigs()
	if err != nil {
		return nil, err
	}
	return dependencyDockerConfigs.Merge(serviceDockerConfigs), nil
}

func (i *Initializer) getServiceDockerConfigs() (dockercompose.DockerConfigs, error) {
	result := dockercompose.DockerConfigs{}
	for serviceName, serviceConfig := range i.ServiceConfigs {
		dockerComposeBuilder := dockercomposebuilder.New(i.AppConfig, serviceConfig, i.ServiceData[serviceName], serviceName, i.AppDir, i.HomeDir)
		dockerConfig, err := dockerComposeBuilder.GetServiceDockerConfigs()
		if err != nil {
			return result, err
		}
		result = result.Merge(result, dockerConfig)
	}
	return result, nil
}

func (i *Initializer) renderDockerCompose(dockerComposeDir string) error {
	bytes, err := yaml.Marshal(i.DockerComposeConfig)
	if err != nil {
		return err
	}
	if err := osplus.CreateEmptyDirectory(dockerComposeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dockerComposeDir, "docker-compose.yml"), bytes, 0777)
}

func (i *Initializer) setupDockerImages(dockerComposeDir string) error {
	if err := dockercompose.PullAllImages(dockerComposeDir, i.logChannel); err != nil {
		return err
	}
	return dockercompose.BuildAllImages(dockerComposeDir, i.logChannel)
}

// Initialize sets up the entire app and returns an error if any
func (i *Initializer) Initialize() error {
	dockerConfigs, err := i.GetDockerConfigs()
	if err != nil {
		return err
	}
	i.DockerComposeConfig.Services = dockerConfigs
	dockerComposeDir := path.Join(i.AppDir, "tmp")
	if err := i.renderDockerCompose(dockerComposeDir); err != nil {
		return err
	}
	return i.setupDockerImages(dockerComposeDir)
}
