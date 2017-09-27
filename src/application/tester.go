package application

import (
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
)

// Tester runs tests for all internal services of the application
type Tester struct {
	AppConfig                types.AppConfig
	InternalServiceConfigs   map[string]types.ServiceConfig
	ServiceData              map[string]types.ServiceData
	AppDir                   string
	homeDir                  string
	DockerComposeProjectName string
	Logger                   *Logger
	logChannel               chan string
	logRole                  string
	mode                     composebuilder.BuildMode
}

// NewTester is Tester's constructor
func NewTester(appConfig types.AppConfig, logger *Logger, appDir, homeDir, dockerComposeProjectName string, mode composebuilder.BuildMode) (*Tester, error) {
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
		Logger:     logger,
		logRole:    "exo-test",
		logChannel: logger.GetLogChannel("exo-test"),
		mode:       mode,
	}, nil
}

// RunAppTests runs the tests for the entire application
func (a *Tester) RunAppTests() (bool, error) {
	a.logChannel <- fmt.Sprintf("Testing application %s", a.AppConfig.Name)
	numFailed := 0
	for serviceName, serviceConfig := range a.InternalServiceConfigs {
		if serviceConfig.Development.Scripts["test"] == "" {
			a.logChannel <- fmt.Sprintf("%s has no tests, skipping", serviceName)
		} else {
			testPassed, err := a.runServiceTests(serviceName, serviceConfig)
			if err != nil {
				a.logChannel <- fmt.Sprintf("error running '%s' tests:", err)
			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		a.logChannel <- "All tests passed"
	} else {
		a.logChannel <- fmt.Sprintf("%d tests failed", numFailed)
	}
	close(a.logChannel)
	a.Logger.WaitForChannelsToClose()
	return numFailed == 0, nil
}

// RunServiceTest runs the tests for a single service
func (a *Tester) RunServiceTest(serviceName string) (bool, error) {
	testsPassed := true
	var err error
	if a.InternalServiceConfigs[serviceName].Development.Scripts["test"] == "" {
		a.logChannel <- fmt.Sprintf("%s has no tests, skipping", serviceName)
	} else {
		if testsPassed, err = a.runServiceTests(serviceName, a.InternalServiceConfigs[serviceName]); err != nil {
			a.logChannel <- fmt.Sprintf("error running '%s' tests:", err)
		}
	}
	return testsPassed, nil
}

// runServiceTests runs the tests for the given service
func (a *Tester) runServiceTests(serviceName string, serviceConfig types.ServiceConfig) (bool, error) {
	a.logChannel <- fmt.Sprintf("Testing service '%s'", serviceName)
	builtDependencies := config.GetBuiltServiceDevelopmentDependencies(serviceConfig, a.AppConfig, a.AppDir, a.homeDir)
	initializer, err := NewInitializer(a.AppConfig, a.logChannel, a.logRole, a.AppDir, a.homeDir, a.DockerComposeProjectName, a.mode)
	if err != nil {
		return false, err
	}
	runner, err := NewRunner(a.AppConfig, a.logChannel, a.AppDir, a.homeDir, a.DockerComposeProjectName)
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
