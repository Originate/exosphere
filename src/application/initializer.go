package application

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/src/config"
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
	ServiceData              map[string]types.ServiceData
	ServiceConfigs           map[string]types.ServiceConfig
	AppDir                   string
	HomeDir                  string
	BuildMode                composebuilder.BuildMode
	logger                   *util.Logger
}

// NewInitializer is Initializer's constructor
func NewInitializer(appConfig types.AppConfig, logger *util.Logger, appDir, homeDir, dockerComposeProjectName string, mode composebuilder.BuildMode) (*Initializer, error) {
	serviceConfigs, err := config.GetServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Initializer{}, err
	}
	appSetup := &Initializer{
		AppConfig:                appConfig,
		DockerComposeConfig:      types.DockerCompose{Version: "3"},
		DockerComposeProjectName: dockerComposeProjectName,
		ServiceData:              appConfig.GetServiceData(),
		ServiceConfigs:           serviceConfigs,
		AppDir:                   appDir,
		HomeDir:                  homeDir,
		BuildMode:                mode,
		logger:                   logger,
	}
	return appSetup, nil
}

func (i *Initializer) getAppDependenciesDockerConfigs() (types.DockerConfigs, error) {
	result := types.DockerConfigs{}
	if i.BuildMode == composebuilder.BuildModeDeployProduction {
		appDependencies := config.GetBuiltAppProductionDependencies(i.AppConfig, i.AppDir)
		for _, builtDependency := range appDependencies {
			if builtDependency.HasDockerConfig() {
				dockerConfig, err := builtDependency.GetDockerConfig()
				if err != nil {
					return result, err
				}
				result[builtDependency.GetServiceName()] = dockerConfig
			}
		}
	} else {
		appDependencies := config.GetBuiltAppDevelopmentDependencies(i.AppConfig, i.AppDir, i.HomeDir)
		for _, builtDependency := range appDependencies {
			dockerConfig, err := builtDependency.GetDockerConfig()
			if err != nil {
				return result, err
			}
			result[builtDependency.GetContainerName()] = dockerConfig
		}
	}
	return result, nil
}

// GetDockerConfigs returns the docker configs of all services and dependencies in
// the application
func (i *Initializer) GetDockerConfigs() (types.DockerConfigs, error) {
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
	for serviceRole, serviceConfig := range i.ServiceConfigs {
		dockerConfig, err := composebuilder.GetServiceDockerConfigs(i.AppConfig, serviceConfig, i.ServiceData[serviceRole], serviceRole, i.AppDir, i.HomeDir, i.BuildMode)
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
