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
	AppConfig           types.AppConfig
	Logger              *logger.Logger
	DockerComposeConfig types.DockerCompose
	ServiceData         map[string]types.ServiceData
	ServiceConfigs      map[string]types.ServiceConfig
	AppDir              string
	HomeDir             string
}

// NewAppSetup is AppSetup's constructor
func NewAppSetup(appConfig types.AppConfig, logger *logger.Logger, appDir, homeDir string) (*AppSetup, error) {
	serviceConfigs, err := serviceConfigHelpers.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &AppSetup{}, err
	}
	appSetup := &AppSetup{
		AppConfig:           appConfig,
		Logger:              logger,
		DockerComposeConfig: types.DockerCompose{Version: "3"},
		ServiceData:         serviceConfigHelpers.GetServiceData(appConfig.Services),
		ServiceConfigs:      serviceConfigs,
		AppDir:              appDir,
		HomeDir:             homeDir,
	}
	return appSetup, nil
}

func (a *AppSetup) getAppDependenciesDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for _, dependency := range a.AppConfig.Dependencies {
		builtDependency := appDependencyHelpers.Build(dependency, a.AppConfig, a.AppDir, a.HomeDir)
		dockerConfig, err := builtDependency.GetDockerConfig()
		if err != nil {
			return result, err
		}
		result[builtDependency.GetContainerName()] = dockerConfig
	}
	return result, nil
}

func (a *AppSetup) getDockerConfigs() (map[string]types.DockerConfig, error) {
	dependencyDockerConfigs, err := a.getAppDependenciesDockerConfigs()
	if err != nil {
		return nil, err
	}
	serviceDockerConfigs, err := a.getServiceDockerConfigs()
	if err != nil {
		return nil, err
	}
	return util.JoinDockerConfigMaps(dependencyDockerConfigs, serviceDockerConfigs), nil
}

func (a *AppSetup) getServiceDockerConfigs() (map[string]types.DockerConfig, error) {
	result := map[string]types.DockerConfig{}
	for serviceName, serviceConfig := range a.ServiceConfigs {
		setup := &dockerSetup.DockerSetup{
			AppConfig:     a.AppConfig,
			ServiceConfig: serviceConfig,
			ServiceData:   a.ServiceData[serviceName],
			Role:          serviceName,
			Logger:        a.Logger,
			AppDir:        a.AppDir,
			HomeDir:       a.HomeDir,
		}
		dockerConfig, err := setup.GetServiceDockerConfigs()
		if err != nil {
			return result, err
		}
		result = util.JoinDockerConfigMaps(result, dockerConfig)
	}
	return result, nil
}

func (a *AppSetup) renderDockerCompose(dockerComposeDir string) error {
	bytes, err := yaml.Marshal(a.DockerComposeConfig)
	if err != nil {
		return err
	}
	if err := osHelpers.EmptyDir(dockerComposeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dockerComposeDir, "docker-compose.yml"), bytes, 0777)
}

func (a *AppSetup) setupDockerImages(dockerComposeDir string) error {
	if err := dockerCompose.PullAllImages(dockerComposeDir, a.Write); err != nil {
		return err
	}
	return dockerCompose.BuildAllImages(dockerComposeDir, a.Write)
}

// Setup sets up the entire app and returns an error if any
func (a *AppSetup) Setup() error {
	dockerConfigs, err := a.getDockerConfigs()
	if err != nil {
		return err
	}
	a.DockerComposeConfig.Services = dockerConfigs
	dockerComposeDir := path.Join(a.AppDir, "tmp")
	if err := a.renderDockerCompose(dockerComposeDir); err != nil {
		return err
	}
	return a.setupDockerImages(dockerComposeDir)
}

// Write logs exo-run output
func (a *AppSetup) Write(text string) {
	err := a.Logger.Log("exo-run", text, true)
	if err != nil {
		fmt.Printf("Error logging exo-run output: %v\n", err)
	}
}
