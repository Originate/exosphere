package application

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// ServiceTester runs the tests for the given service
type ServiceTester struct {
	DirName           string
	ServiceConfig     types.ServiceConfig
	BuiltDependencies map[string]config.AppDevelopmentDependency
	AppDir            string
	ServiceLocation   string
	Initializer       *Initializer
	Runner            *Runner
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

	initializer, err := NewInitializer(
		appConfig,
		logger,
		appDir,
		homeDir,
		dockerComposeProjectName,
		mode,
	)
	if err != nil {
		return nil, err
	}

	builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfig, appConfig, appDir, homeDir)

	runner, err := NewRunner(appConfig, logger, appDir, homeDir, dockerComposeProjectName)

	if err != nil {
		return nil, err
	}
	tester := &ServiceTester{
		DirName:           serviceDir,
		ServiceConfig:     serviceConfig,
		BuiltDependencies: builtDependencies,
		AppDir:            appDir,
		ServiceLocation:   serviceLocation,
		Initializer:       initializer,
		Runner:            runner,
	}
	tester.Initializer.DockerComposeConfig, err = tester.getDockerComposeConfig()
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

func (s *ServiceTester) getDockerComposeConfig() (types.DockerCompose, error) {
	result := types.DockerCompose{Version: "3"}
	dependencyDockerConfigs, err := composebuilder.GetDependenciesDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: s.Initializer.AppConfig,
		AppDir:    s.AppDir,
		BuildMode: s.Initializer.BuildMode,
		HomeDir:   s.Initializer.HomeDir,
	})
	if err != nil {
		return result, err
	}
	dockerConfigs, err := composebuilder.GetServiceDockerConfigs(
		s.Initializer.AppConfig,
		s.ServiceConfig,
		types.ServiceData{Location: s.DirName},
		s.DirName,
		s.AppDir,
		s.Initializer.HomeDir,
		s.Initializer.BuildMode,
	)
	if err != nil {
		return result, err
	}
	serviceConfig := dockerConfigs[s.DirName]
	serviceConfig.DependsOn = s.getDependencyContainerNames()
	serviceConfig.Command = s.ServiceConfig.Development.Scripts["test"]
	dockerConfigs[s.DirName] = serviceConfig
	for _, builtDependency := range s.BuiltDependencies {
		dockerConfigs[builtDependency.GetContainerName()] = dependencyDockerConfigs[builtDependency.GetContainerName()]
	}
	result.Services = dockerConfigs
	return result, nil
}

func (s *ServiceTester) getServiceNames() []string {
	return []string{s.DirName}
}

func (s *ServiceTester) getServiceOnlineTexts() map[string]string {
	return map[string]string{
		"": fmt.Sprintf("%s exited with code", s.DirName),
	}
}

func (s *ServiceTester) runTests() (int, error) {
	dependencyNames := s.getDependencyContainerNames()
	if len(dependencyNames) > 0 {
		if _, err := s.Runner.runImages(dependencyNames, s.getDependencyOnlineTexts(), "dependencies"); err != nil {
			return 1, err
		}
	}
	output, err := s.Runner.runImages(s.getServiceNames(), s.getServiceOnlineTexts(), "services")
	if err != nil {
		return 1, err
	}
	return util.GetServiceExitCode(s.DirName, output)
}

func (s *ServiceTester) setup() error {
	dockerComposeDir := path.Join(s.ServiceLocation, "tests", "tmp")
	if err := s.Initializer.renderDockerCompose(dockerComposeDir); err != nil {
		return err
	}
	if err := s.Initializer.setupDockerImages(dockerComposeDir); err != nil {
		return err
	}
	s.Runner.DockerComposeDir = dockerComposeDir
	return nil
}

// Start sets up and runs the test for a service, return its result and an error if any
func (s *ServiceTester) Start() (int, error) {
	if err := s.setup(); err != nil {
		return 1, err
	}

	exitCode, err := s.runTests()
	return exitCode, err
}

// Shutdown shuts down the tests
func (s *ServiceTester) Shutdown() error {
	return s.Runner.Shutdown(types.ShutdownConfig{CloseMessage: "killing test containers\n"})
}
