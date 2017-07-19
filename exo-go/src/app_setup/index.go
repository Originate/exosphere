package appSetup

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/Originate/exosphere/exo-go/src/docker_setup"
	"github.com/Originate/exosphere/exo-go/src/logger"
	"github.com/Originate/exosphere/exo-go/src/service_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// AppSetup sets up the app
type AppSetup struct {
	AppConfig             types.AppConfig
	Logger                *logger.Logger
	DockerComposeConfig   types.DockerCompose
	ServiceData           map[string]types.ServiceData
	ServiceConfigs        map[string]types.ServiceConfig
	DockerComposeLocation string
	Cwd                   string
}

// NewAppSetup is AppSetup's constructor
func NewAppSetup(appConfig types.AppConfig, logger *logger.Logger) (*AppSetup, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(cwd, appConfig)
	if err != nil {
		return &AppSetup{}, err
	}
	appSetup := &AppSetup{
		AppConfig:             appConfig,
		Logger:                logger,
		DockerComposeConfig:   types.DockerCompose{Version: "3"},
		DockerComposeLocation: path.Join(cwd, "tmp"),
		ServiceConfigs:        serviceConfigs,
		Cwd:                   cwd,
	}
	return appSetup, nil
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
	return joinDockerConfigMaps(dependencyDockerConfigs, serviceDockerConfigs), nil
}

func (appSetup *AppSetup) getAppDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range appSetup.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, appSetup.AppConfig)
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

func (appSetup *AppSetup) getServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for serviceName, serviceConfig := range appSetup.ServiceConfigs {
		dockerConfig, err := dockerSetup.NewDockerSetup(appSetup.AppConfig, serviceConfig, appSetup.ServiceData[serviceName], serviceName, appSetup.Logger).GetServiceDockerConfig()
		if err != nil {
			return result, err
		}
		result = joinDockerConfigMaps(result, dockerConfig)
	}
	return result, nil
}

func (appSetup *AppSetup) renderDockerCompose() error {
	bytes, err := yaml.Marshal(appSetup.DockerComposeConfig)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(appSetup.DockerComposeLocation, bytes, 0777)
}

func (appSetup *AppSetup) setupDockerImages() error {
	if _, err := dockerCompose.PullAllImages(appSetup.Cwd, appSetup.Write); err != nil {
		return err
	}
	if _, err := dockerCompose.BuildAllImages(appSetup.Cwd, appSetup.Write); err != nil {
		appSetup.Write("Docker setup failed")
		return err
	}
	appSetup.Write("Docker setup finished")
	return nil
}

// StartSetup sets up the entire app
func (appSetup *AppSetup) StartSetup() error {
	dockerConfigs, err := appSetup.getDockerConfigs()
	if err != nil {
		return err
	}
	for service, dockerConfig := range dockerConfigs {
		appSetup.DockerComposeConfig.Services[service] = dockerConfig
	}
	return nil
}

// Write logs exo-run output
func (appSetup *AppSetup) Write(text string) {
	appSetup.Logger.Log("exo-run", text, true)
}
