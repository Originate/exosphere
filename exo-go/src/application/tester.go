package application

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/types"
)

// Tester runs tests for all internal services of the application
type Tester struct {
	AppConfig              types.AppConfig
	InternalServiceConfigs map[string]types.ServiceConfig
	ServiceData            map[string]types.ServiceData
	AppDir                 string
	homeDir                string
	Logger                 *Logger
	logChannel             chan string
	logRole                string
}

// NewTester is Tester's constructor
func NewTester(appConfig types.AppConfig, logger *Logger, appDir, homeDir string) (*Tester, error) {
	internalServiceConfigs, err := config.GetInternalServiceConfigs(appDir, appConfig)
	if err != nil {
		return &Tester{}, err
	}
	return &Tester{
		AppConfig:              appConfig,
		InternalServiceConfigs: internalServiceConfigs,
		ServiceData:            config.GetServiceData(appConfig.Services),
		AppDir:                 appDir,
		homeDir:                homeDir,
		Logger:                 logger,
		logRole:                "exo-test",
		logChannel:             logger.GetLogChannel("exo-test"),
	}, nil
}

// RunAppTests runs the tests for the entire application
func (a *Tester) RunAppTests(serviceName string) error {
	a.logChannel <- fmt.Sprintf("Testing application %s", a.AppConfig.Name)
	numFailed := 0
	for serviceName, serviceConfig := range a.InternalServiceConfigs {
		if serviceConfig.Tests == "" {
			a.logChannel <- fmt.Sprintf("%s has no tests, skipping", serviceName)
		} else {
			if testPassed, err := a.runServiceTests(serviceName, serviceConfig); err != nil {
				a.logChannel <- fmt.Sprintf("error running '%s' tests:", err)
			} else if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		return a.Logger.Log("exo-test", "All tests passed", true)
	}
	return a.Logger.Log("exo-test", fmt.Sprintf("%d tests failed", numFailed), true)
}

// runServiceTests runs the tests for the given service
func (a *Tester) runServiceTests(serviceName string, serviceConfig types.ServiceConfig) (bool, error) {
	a.logChannel <- fmt.Sprintf("Testing service '%s'", serviceName)
	builtDependencies := config.GetServiceBuiltDependencies(serviceConfig, a.AppConfig, a.AppDir, a.homeDir)
	initializer, err := NewInitializer(a.AppConfig, a.Logger, a.logRole, a.AppDir, a.homeDir)
	if err != nil {
		return false, err
	}
	runner, err := NewRunner(a.AppConfig, a.Logger, a.logRole, a.AppDir, a.homeDir)
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
