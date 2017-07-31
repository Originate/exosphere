package application

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/docker_setup"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// Initializer sets up the app
type Initializer struct {
	AppConfig           types.AppConfig
	Logger              *Logger
	DockerComposeConfig types.DockerCompose
	ServiceData         map[string]types.ServiceData
	ServiceConfigs      map[string]types.ServiceConfig
	AppDir              string
	HomeDir             string
	logChannel          chan string
}

// NewInitializer is Initializer's constructor
func NewInitializer(appConfig types.AppConfig, logger *Logger, appDir, homeDir string) (*Initializer, error) {
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Initializer{}, err
	}
	appSetup := &Initializer{
		AppConfig:           appConfig,
		Logger:              logger,
		DockerComposeConfig: types.DockerCompose{Version: "3"},
		ServiceData:         serviceConfigHelpers.GetServiceData(appConfig.Services),
		ServiceConfigs:      serviceConfigs,
		AppDir:              appDir,
		HomeDir:             homeDir,
		logChannel:          logger.GetLogChannel("exo-run"),
	}
	return appSetup, nil
}

func (i *Initializer) getAppDependenciesDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	for _, dependency := range i.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, i.AppConfig, i.AppDir, i.HomeDir)
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

func (i *Initializer) getDockerConfigs() (types.DockerConfigs, error) {
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

func (i *Initializer) getServiceDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	for serviceName, serviceConfig := range i.ServiceConfigs {
		setup := &dockerSetup.DockerSetup{
			AppConfig:     i.AppConfig,
			ServiceConfig: serviceConfig,
			ServiceData:   i.ServiceData[serviceName],
			Role:          serviceName,
			AppDir:        i.AppDir,
			HomeDir:       i.HomeDir,
		}
		dockerConfig, err := setup.GetServiceDockerConfigs()
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
	if err := util.EmptyDir(dockerComposeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dockerComposeDir, "docker-compose.yml"), bytes, 0777)
}

func (i *Initializer) setupDockerImages(dockerComposeDir string) error {
	if err := dockerCompose.PullAllImages(dockerComposeDir, i.logChannel); err != nil {
		return err
	}
	return dockerCompose.BuildAllImages(dockerComposeDir, i.logChannel)
}

// Initialize sets up the entire app and returns an error if any
func (i *Initializer) Initialize() error {
	dockerConfigs, err := i.getDockerConfigs()
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
