package application

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// Tester runs tests for all internal services of the application
type Tester struct {
	AppConfig                types.AppConfig
	InternalServiceConfigs   map[string]types.ServiceConfig
	ServiceData              map[string]types.ServiceData
	AppDir                   string
	homeDir                  string
	DockerComposeProjectName string
	logger                   *util.Logger
	mode                     composebuilder.BuildMode
}

// NewTester is Tester's constructor
func NewTester(appConfig types.AppConfig, logger *util.Logger, appDir, homeDir, dockerComposeProjectName string, mode composebuilder.BuildMode) (*Tester, error) {
	internalServiceConfigs, err := config.GetInternalServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Tester{}, err
	}
	return &Tester{
		AppConfig:                appConfig,
		InternalServiceConfigs:   internalServiceConfigs,
		ServiceData:              appConfig.GetServiceData(),
		AppDir:                   appDir,
		homeDir:                  homeDir,
		DockerComposeProjectName: dockerComposeProjectName,
		logger: logger,
		mode:   mode,
	}, nil
}

// RunAppTests runs the tests for the entire application
func (a *Tester) RunAppTests() (bool, error) {
	a.logger.Logf("Testing application %s", a.AppConfig.Name)
	numFailed := 0
	for serviceName, serviceConfig := range a.InternalServiceConfigs {
		if serviceConfig.Development.Scripts["test"] == "" {
			a.logger.Logf("%s has no tests, skipping", serviceName)
		} else {
			testPassed, err := a.runServiceTests(serviceName, serviceConfig)
			if err != nil {
				a.logger.Logf("error running '%s' tests:", err)
			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		a.logger.Log("All tests passed")
	} else {
		a.logger.Logf("%d tests failed", numFailed)
	}
	return numFailed == 0, nil
}

// RunServiceTest runs the tests for a single service
func (a *Tester) RunServiceTest(serviceName string) (bool, error) {
	testsPassed := true
	var err error
	if a.InternalServiceConfigs[serviceName].Development.Scripts["test"] == "" {
		a.logger.Logf("%s has no tests, skipping", serviceName)
	} else {
		if testsPassed, err = a.runServiceTests(serviceName, a.InternalServiceConfigs[serviceName]); err != nil {
			a.logger.Logf("error running '%s' tests:", err)
		}
	}
	return testsPassed, nil
}

// runServiceTests runs the tests for the given service
func (a *Tester) runServiceTests(serviceName string, serviceConfig types.ServiceConfig) (bool, error) {
	a.logger.Logf("Testing service '%s'", serviceName)
	builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfig, a.AppConfig, a.AppDir, a.homeDir)
	initializer, err := NewInitializer(a.AppConfig, a.logger, a.AppDir, a.homeDir, a.DockerComposeProjectName, a.mode)
	if err != nil {
		return false, err
	}
	runner, err := NewRunner(a.AppConfig, a.logger, a.AppDir, a.homeDir, a.DockerComposeProjectName)
	if err != nil {
		return false, err
	}
	if err != nil {
		return false, err
	}
	serviceTester, err := NewServiceTester(serviceName, serviceConfig, builtDependencies, a.AppDir, a.ServiceData[serviceName].Location, initializer, runner)
	if err != nil {
		return false, err
	}
	return serviceTester.Run()
}
