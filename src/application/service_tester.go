package application

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/runner"
	"github.com/Originate/exosphere/src/docker/tools"
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
	Logger                   *util.Logger
	RunOptions               runner.RunOptions
}

// NewServiceTester is ServiceTester's constructor
func NewServiceTester(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode) (*ServiceTester, error) {
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
		Logger:                   logger,
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

func (s *ServiceTester) getDependencyOnlineTexts() map[string]string {
	result := map[string]string{}
	for dependencyName, builtDependency := range s.BuiltDependencies {
		result[dependencyName] = builtDependency.GetOnlineText()
	}
	return result
}

func (s *ServiceTester) getDockerComposeConfig() (types.DockerConfigs, error) {
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
		return dockerConfigs, err
	}
	dependencyDockerConfigs, err := composebuilder.GetDependenciesDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: s.AppConfig,
		AppDir:    s.AppDir,
		BuildMode: s.BuildMode,
		HomeDir:   s.HomeDir,
	})
	if err != nil {
		return dockerConfigs, err
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

func (s *ServiceTester) getServiceNames() []string {
	return []string{s.DirName}
}

func (s *ServiceTester) getServiceOnlineTexts() map[string]string {
	return map[string]string{
		"": fmt.Sprintf("%s exited with code", s.DirName),
	}
}

// Run runs the tests for the service and return true if the tests passed
// and an error if any
func (s *ServiceTester) Run() (int, error) {
	err := runner.Run(s.RunOptions)
	if err != nil {
		return 1, err
	}
	return tools.GetExitCode(s.DirName)
}

// Shutdown shuts down the tests
func (s *ServiceTester) Shutdown() error {
	return runner.Shutdown(s.RunOptions)
}

func (s *ServiceTester) getRunOptions() (runner.RunOptions, error) {
	dockerComposeDir := path.Join(s.ServiceLocation, "tests", "tmp")
	dockerConfigs, err := s.getDockerComposeConfig()
	if err != nil {
		return runner.RunOptions{}, err
	}
	return runner.RunOptions{
		ImageGroups: []runner.ImageGroup{
			{
				ID:          "dependencies",
				Names:       s.getDependencyContainerNames(),
				OnlineTexts: s.getDependencyOnlineTexts(),
			},
			{
				ID:          "services",
				Names:       s.getServiceNames(),
				OnlineTexts: s.getServiceOnlineTexts(),
			},
		},
		DockerConfigs:            dockerConfigs,
		DockerComposeDir:         dockerComposeDir,
		DockerComposeProjectName: s.DockerComposeProjectName,
		Logger: s.Logger,
	}, nil
}
