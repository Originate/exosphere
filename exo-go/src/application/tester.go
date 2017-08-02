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
		logChannel:             logger.GetLogChannel("exo-test"),
	}, nil
}

// Run runs the tests
func (a *Tester) Run() {
	testErrChannel := make(chan error)
	numFinished, numFailed := 0, 0
	go func() {
		for numFinished != len(a.InternalServiceConfigs) {
			if err := <-testErrChannel; err != nil {
				numFailed++
			}
			numFinished++
		}
		if numFailed == 0 {
			a.logChannel <- "All tests passed"
			return
		}
		a.logChannel <- fmt.Sprintf("%d tests failed", numFailed)
		return
	}()
	for serviceName, serviceConfig := range a.InternalServiceConfigs {
		if serviceConfig.Tests == "" {
			a.logChannel <- fmt.Sprintf("%s has no tests, skipping", serviceName)
		} else {
			a.runServiceTests(serviceName, serviceConfig, testErrChannel)
		}
	}
}

func (a *Tester) runServiceTests(serviceName string, serviceConfig types.ServiceConfig, testErrChannel chan error) {
	builtDependencies := config.GetServiceBuiltDependencies(serviceConfig, a.AppConfig, a.AppDir, a.homeDir)
	initializer, err := NewInitializer(a.AppConfig, a.Logger, a.AppDir, a.homeDir)
	initializer.logChannel = a.logChannel
	if err != nil {
		testErrChannel <- err
	}
	runner, err := NewRunner(a.AppConfig, a.Logger, a.AppDir, a.homeDir)
	if err != nil {
		testErrChannel <- err
	}
	runner.logChannel = a.logChannel
	if err != nil {
		testErrChannel <- err
	}
	serviceTester, err := NewServiceTester(serviceName, serviceConfig, builtDependencies, a.AppDir, a.ServiceData[serviceName].Location, initializer, runner)
	if err != nil {
		testErrChannel <- err
	}
	serviceTester.Run(testErrChannel)
}
