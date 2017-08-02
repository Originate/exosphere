package application

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/util"
)

// ServiceTester runs the tests for the given service
type ServiceTester struct {
	Role              string
	ServiceConfig     types.ServiceConfig
	BuiltDependencies map[string]config.AppDependency
	AppDir            string
	ServiceDir        string
	*Initializer
	*Runner
}

// NewServiceTester is ServiceTester's constructor
func NewServiceTester(role string, serviceConfig types.ServiceConfig, builtDependencies map[string]config.AppDependency, appDir, serviceDir string, initializer *Initializer, runner *Runner) (*ServiceTester, error) {
	tester := &ServiceTester{
		Role:              role,
		ServiceConfig:     serviceConfig,
		BuiltDependencies: builtDependencies,
		AppDir:            appDir,
		ServiceDir:        serviceDir,
		Initializer:       initializer,
		Runner:            runner,
	}
	var err error
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
	dockerConfigs := types.DockerConfigs{}
	serviceDockerConfig := appDockerConfigs[s.Role]
	serviceDockerConfig.Build = map[string]string{
		"context":    "../../",
		"dockerfile": "tests/Dockerfile",
	}
	serviceDockerConfig.Command = s.ServiceConfig.Tests
	dockerConfigs[s.Role] = serviceDockerConfig
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
	if _, err := s.runImages(s.getDependencyContainerNames(), s.getDependencyOnlineTexts(), "dependencies"); err != nil {
		return 1, err
	}
	output, err := s.runImages(s.getServiceNames(), s.getServiceOnlineTexts(), "services")
	if err != nil {
		return 1, err
	}
	return util.GetServiceExitCode(s.Role, output)
}

func (s *ServiceTester) setup() error {
	dockerComposeDir := path.Join(s.AppDir, s.ServiceDir, "tests", "tmp")
	if err := s.renderDockerCompose(dockerComposeDir); err != nil {
		return err
	}
	if err := s.setupDockerImages(dockerComposeDir); err != nil {
		return err
	}
	s.Runner.DockerComposeDir = dockerComposeDir
	return nil
}

// Run runs runs the tests for the service
func (s *ServiceTester) Run(testErrChannel chan error) {
	if err := s.setup(); err != nil {
		testErrChannel <- err
		return
	}
	exitCode, err := s.runTests()
	if err != nil {
		testErrChannel <- err
		return
	}
	if err := s.Shutdown(types.ShutdownConfig{CloseMessage: "killing test containers\n"}); err != nil {
		testErrChannel <- err
		return
	}
	if exitCode > 0 {
		s.Runner.logChannel <- fmt.Sprintf("'%s' tests failed", s.Role)
		testErrChannel <- fmt.Errorf("'%s' tests failed", s.Role)
	} else {
		s.Runner.logChannel <- fmt.Sprintf("'%s' tests passed", s.Role)
		testErrChannel <- nil
	}
}
