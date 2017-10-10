package application

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, logger *util.Logger, mode composebuilder.BuildMode) (bool, error) {
	logger.Logf("Testing application %s", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return false, err
	}

	numFailed := 0
	for _, serviceContext := range serviceContexts {
		if serviceContext.Config.Development.Scripts["test"] == "" {
			logger.Logf("%s has no tests, skipping", serviceContext.Dir)
		} else {
			testPassed, err := TestService(serviceContext, logger, mode)
			if err != nil {
				logger.Logf("error running '%s' tests:", err)
			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		logger.Log("All tests passed")
		return true, nil
	}
	logger.Logf("%d tests failed", numFailed)
	return false, nil
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode) (bool, error) {
	logger.Logf("Testing service '%s'", serviceContext.Dir)
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
	logger.Logf("'%s' tests %s", serviceContext.Dir, result)
	return testPassed, serviceTester.Shutdown()
}
