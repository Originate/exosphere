package application

import (
	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed, if they were interrupted,
// and an error if any
func TestApp(appContext types.AppContext, logger *util.Logger, mode composebuilder.BuildMode) types.TestResult {
	logger.Logf("Testing application %s", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return types.TestResult{
			Passed:      false,
			Interrupted: false,
			Error:       err,
		}
	}

	numFailed := 0
	locations := []string{}
	for _, serviceContext := range serviceContexts {
		if util.DoesStringArrayContain(locations, serviceContext.Location) {
			continue
		}
		locations = append(locations, serviceContext.Location)
		if serviceContext.Config.Development.Scripts["test"] == "" {
			logger.Logf("%s has no tests, skipping", serviceContext.Dir)
		} else {
			testResult := TestService(serviceContext, logger, mode)
			switch {
			case testResult.Interrupted:
				return testResult
			case testResult.Error != nil:
				logger.Logf("error running '%s' tests:", testResult.Error)
			case testResult.Passed:
				logger.Logf("'%s' tests passed", serviceContext.Dir)
			case !testResult.Passed:
				logger.Logf("'%s' tests failed", serviceContext.Dir)
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		logger.Log("All tests passed")
		return types.TestResult{
			Passed:      true,
			Interrupted: false,
			Error:       nil,
		}
	}
	logger.Logf("%d tests failed", numFailed)
	return types.TestResult{
		Passed:      false,
		Interrupted: false,
		Error:       nil,
	}
}

// TestService runs the tests for the service and returns a TestResult struct
func TestService(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode) types.TestResult {
	logger.Logf("Testing service '%s'", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, logger, mode)
	if err != nil {
		return types.TestResult{
			Passed:      false,
			Interrupted: false,
			Error:       err,
		}
	}

	return serviceTester.Run()
}
