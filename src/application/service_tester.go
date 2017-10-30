package application

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// ServiceTester runs the tests for the given service
type ServiceTester struct {
	AppConfig                types.AppConfig
	HomeDir                  string
	DirName                  string
	ServiceConfig            types.ServiceConfig
	BuiltDependencies        map[string]config.AppDevelopmentDependency
	AppDir                   string
	ServiceLocation          string
	BuildMode                composebuilder.BuildMode
	DockerComposeProjectName string
	Writer                   io.Writer
	RunOptions               composerunner.RunOptions
}

// NewServiceTester is ServiceTester's constructor
func NewServiceTester(serviceContext types.ServiceContext, writer io.Writer, mode composebuilder.BuildMode) (*ServiceTester, error) {
	appConfig := serviceContext.AppContext.Config
	serviceConfig := serviceContext.Config
	serviceLocation := serviceContext.Location
	serviceDir := serviceContext.Dir
	appDir := serviceContext.AppContext.Location
	dockerComposeProjectName := fmt.Sprintf("%stests", composebuilder.GetDockerComposeProjectName(appDir))
	homeDir, err := util.GetHomeDirectory()
	if err != nil {
		return nil, err
	}
	builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir, homeDir)
	if err != nil {
		return nil, err
	}
	tester := &ServiceTester{
		AppConfig:                appConfig,
		AppDir:                   appDir,
		BuildMode:                mode,
		BuiltDependencies:        builtDependencies,
		DirName:                  serviceDir,
		DockerComposeProjectName: dockerComposeProjectName,
		HomeDir:                  homeDir,
		Writer:                   writer,
		ServiceConfig:            serviceConfig,
		ServiceLocation:          serviceLocation,
	}
	tester.RunOptions, err = tester.getRunOptions()
	if err != nil {
		return nil, err
	}
	return tester, err
}

func (s *ServiceTester) getDependencyContainerNames() []string {
	dependencyNames := []string{}
	for _, builtDependency := range s.BuiltDependencies {
		dependencyNames = append(dependencyNames, builtDependency.GetContainerName())
	}
	return dependencyNames
}

func (s *ServiceTester) getDockerComposeConfig() (types.DockerConfigs, error) {
	dependencyDockerConfigs, err := composebuilder.GetDependenciesDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: s.AppConfig,
		AppDir:    s.AppDir,
		BuildMode: s.BuildMode,
		HomeDir:   s.HomeDir,
	})
	if err != nil {
		return types.DockerConfigs{}, err
	}
	dockerConfigs, err := composebuilder.GetServiceDockerConfigs(
		s.AppConfig,
		s.ServiceConfig,
		types.ServiceData{Location: s.DirName},
		s.DirName,
		s.AppDir,
		s.HomeDir,
		s.BuildMode,
	)
	if err != nil {
		return types.DockerConfigs{}, err
	}
	serviceConfig := dockerConfigs[s.DirName]
	serviceConfig.DependsOn = s.getDependencyContainerNames()
	serviceConfig.Command = s.ServiceConfig.Development.Scripts["test"]
	dockerConfigs[s.DirName] = serviceConfig
	for _, builtDependency := range s.BuiltDependencies {
		dockerConfigs[builtDependency.GetContainerName()] = dependencyDockerConfigs[builtDependency.GetContainerName()]
	}
	return dockerConfigs, nil
}

// Run runs the tests for the service and return true if the tests passed
// and an error if any
func (s *ServiceTester) Run() (int, error) {
	err := composerunner.Run(s.RunOptions)
	if err != nil {
		if strings.Contains(err.Error(), "exit status") {
			return 1, nil
		}
		return 1, err
	}
	return 0, nil
}

// Shutdown shuts down the tests
func (s *ServiceTester) Shutdown() error {
	return composerunner.Shutdown(s.RunOptions)
}

func (s *ServiceTester) getRunOptions() (composerunner.RunOptions, error) {
	dockerComposeDir := path.Join(s.ServiceLocation, "tests", "tmp")
	dockerConfigs, err := s.getDockerComposeConfig()
	if err != nil {
		return composerunner.RunOptions{}, err
	}
	return composerunner.RunOptions{
		DockerConfigs:            dockerConfigs,
		DockerComposeDir:         dockerComposeDir,
		DockerComposeProjectName: s.DockerComposeProjectName,
		Writer:      s.Writer,
		AbortOnExit: true,
	}, nil
}
