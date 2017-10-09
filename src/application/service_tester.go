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
	Role              string
	ServiceConfig     types.ServiceConfig
	BuiltDependencies map[string]config.AppDevelopmentDependency
	AppDir            string
	ServiceDir        string
	Initializer       *Initializer
	Runner            *Runner
}

// NewServiceTester is ServiceTester's constructor
func NewServiceTester(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode) (*ServiceTester, error) {
	appConfig := serviceContext.AppContext.Config
	serviceConfig := serviceContext.Config
	serviceDir := serviceContext.Location
	role := serviceContext.Name
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
		Role:              role,
		ServiceConfig:     serviceConfig,
		BuiltDependencies: builtDependencies,
		AppDir:            appDir,
		ServiceDir:        serviceDir,
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
	appDockerConfigs, err := s.Initializer.GetDockerConfigs()
	if err != nil {
		return result, err
	}
	dockerConfigs, err := s.Initializer.getServiceDockerConfig(s.Role, s.ServiceConfig, types.ServiceData{Location: s.Role})
	if err != nil {
		return result, err
	}
	serviceConfig := dockerConfigs[s.Role]
	serviceConfig.DependsOn = s.getDependencyContainerNames()
	serviceConfig.Command = s.ServiceConfig.Development.Scripts["test"]
	for _, builtDependency := range s.BuiltDependencies {
		dockerConfigs[builtDependency.GetContainerName()] = appDockerConfigs[builtDependency.GetContainerName()]
	}
	result.Services = dockerConfigs
	return result, nil
}

func (s *ServiceTester) getServiceNames() []string {
	return []string{s.Role}
}

func (s *ServiceTester) getServiceOnlineTexts() map[string]string {
	return map[string]string{
		"": fmt.Sprintf("%s exited with code", s.Role),
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
	return util.GetServiceExitCode(s.Role, output)
}

func (s *ServiceTester) setup() error {
	dockerComposeDir := path.Join(s.ServiceDir, "tests", "tmp")
	if err := s.Initializer.renderDockerCompose(dockerComposeDir); err != nil {
		return err
	}
	if err := s.Initializer.setupDockerImages(dockerComposeDir); err != nil {
		return err
	}
	s.Runner.DockerComposeDir = dockerComposeDir
	return nil
}

// Run runs the tests for the service and return true if the tests passed
// and an error if any
func (s *ServiceTester) Run() (int, error) {
	if err := s.setup(); err != nil {
		return 1, err
	}
	return s.runTests()
}

// Shutdown shuts down the tests
func (s *ServiceTester) Shutdown() error {
	return s.Runner.Shutdown(types.ShutdownConfig{CloseMessage: "killing test containers\n"})
}
