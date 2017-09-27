package application

import (
	"fmt"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, logger *Logger, mode composebuilder.BuildMode) (bool, error) {
	logChannel := logger.GetLogChannel("exo-test")
	logChannel <- fmt.Sprintf("Testing application %s", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return false, err
	}

	numFailed := 0
	for serviceName, serviceContext := range serviceContexts {
		if serviceContext.Config.Development.Scripts["test"] == "" {
			logChannel <- fmt.Sprintf("%s has no tests, skipping", serviceName)
		} else {
			testPassed, err := TestService(serviceContext, logger, mode)
			if err != nil {
				logChannel <- fmt.Sprintf("error running '%s' tests:", err)
			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		logChannel <- "All tests passed"
		return true, nil
	} else {
		logChannel <- fmt.Sprintf("%d tests failed", numFailed)
		return false, nil
	}
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext types.ServiceContext, logger *Logger, mode composebuilder.BuildMode) (bool, error) {
	logChannel := logger.GetLogChannel("exo-test")
	serviceTester, err := NewServiceTester(serviceContext, logger, mode)
	if err != nil {
		return false, err
	}

	exitCode, err := serviceTester.Run()
	if err != nil {
		return false, err
	}
	var testPassed bool
	var result string
	if exitCode == 0 {
		testPassed = true
		result = "passed"
	} else {
		testPassed = false
		result = "failed"
	}
	logChannel <- fmt.Sprintf("'%s' tests %s", serviceContext.Name, result)
	serviceTester.Shutdown()
	return testPassed, nil
}
