package appSetup

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/docker_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// AppSetup sets up the app
type AppSetup struct {
	AppConfig             types.AppConfig
	Logger                *logger.Logger
	DockerComposeConfig   types.DockerCompose
	ServiceData           map[string]types.ServiceData
	ServiceConfigs        map[string]types.ServiceConfig
	DockerComposeLocation string
	AppDir                string
	HomeDir               string
}

// NewAppSetup is AppSetup's constructor
func NewAppSetup(appConfig types.AppConfig, logger *logger.Logger, appDir, homeDir string) (*AppSetup, error) {
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &AppSetup{}, err
	}
	appSetup := &AppSetup{
		AppConfig:             appConfig,
		Logger:                logger,
		DockerComposeConfig:   types.DockerCompose{Version: "3"},
		ServiceData:           serviceConfigHelpers.GetServiceData(appConfig.Services),
		ServiceConfigs:        serviceConfigs,
		DockerComposeLocation: path.Join(appDir, "tmp"),
		AppDir:                appDir,
		HomeDir:               homeDir,
	}
	return appSetup, nil
}

func (appSetup *AppSetup) getAppDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range appSetup.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, appSetup.AppConfig, appSetup.AppDir, appSetup.HomeDir)
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

func (appSetup *AppSetup) getDockerConfigs() (map[string]types.DockerConfig, error) {
	dependencyDockerConfigs, err := appSetup.getAppDependenciesDockerConfigs()
	if err != nil {
		return nil, err
	}
	serviceDockerConfigs, err := appSetup.getServiceDockerConfigs()
	if err != nil {
		return nil, err
	}
	return util.JoinDockerConfigMaps(dependencyDockerConfigs, serviceDockerConfigs), nil
}

func (appSetup *AppSetup) getServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for serviceName, serviceConfig := range appSetup.ServiceConfigs {
		setup := &dockerSetup.DockerSetup{
			AppConfig:     appSetup.AppConfig,
			ServiceConfig: serviceConfig,
			ServiceData:   appSetup.ServiceData[serviceName],
			Role:          serviceName,
			Logger:        appSetup.Logger,
			AppDir:        appSetup.AppDir,
			HomeDir:       appSetup.HomeDir,
		}
		dockerConfig, err := setup.GetServiceDockerConfigs()
		if err != nil {
			return result, err
		}
		result = util.JoinDockerConfigMaps(result, dockerConfig)
	}
	return result, nil
}

func (appSetup *AppSetup) renderDockerCompose() error {
	bytes, err := yaml.Marshal(appSetup.DockerComposeConfig)
	if err != nil {
		return err
	}
	if err := osHelpers.EmptyDir(appSetup.DockerComposeLocation); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(appSetup.DockerComposeLocation, "docker-compose.yml"), bytes, 0777)
}

func (appSetup *AppSetup) setupDockerImages() error {
	if err := dockerCompose.PullAllImages(appSetup.DockerComposeLocation, appSetup.Write); err != nil {
		return err
	}
	return dockerCompose.BuildAllImages(appSetup.DockerComposeLocation, appSetup.Write)
}

// Setup sets up the entire app and returns an error if any
func (appSetup *AppSetup) Setup() error {
	dockerConfigs, err := appSetup.getDockerConfigs()
	if err != nil {
		return err
	}
	appSetup.DockerComposeConfig.Services = dockerConfigs
	if err := appSetup.renderDockerCompose(); err != nil {
		return err
	}
	return appSetup.setupDockerImages()
}

func (appSetup *AppSetup) updateDockerfile() error {
	for serviceName, data := range appSetup.ServiceData {
		if data.Location != "" {
			command := fmt.Sprintf("RUN %s", appSetup.ServiceConfigs[serviceName].Startup["command"])
			osHelpers.AppendToFile(path.Join(data.Location, "Dockerfile"), command)
		}
	}
}

// Write logs exo-run output
func (appSetup *AppSetup) Write(text string) {
	appSetup.Logger.Log("exo-run", text, true)
}
