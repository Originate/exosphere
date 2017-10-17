package application

import (
	"os"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed, if they were interrupted,
// and an error if any
func TestApp(appContext types.AppContext, logger *util.Logger, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	logger.Logf("Testing application %s", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return types.TestResult{}, err
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
			testResult, err := TestService(serviceContext, logger, mode, shutdown)
			switch {
			case testResult.Interrupted:
				return testResult, err
			case err != nil:
				logger.Logf("error running '%s' tests:", err)
				fallthrough
			case !testResult.Passed:
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		logger.Log("All tests passed")
		return types.TestResult{
			Passed: true,
		}, nil
	}
	logger.Logf("%d tests failed", numFailed)
	return types.TestResult{}, nil
}

// TestService runs the tests for the service and returns a TestResult struct
func TestService(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	logger.Logf("Testing service '%s'", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, logger, mode)
	if err != nil {
		return types.TestResult{}, err
	}

	testExit := make(chan int)
	testError := make(chan error)
	go func() {
		exitCode, err := serviceTester.Start()
		if err != nil {
			testError <- err
			return
		}
		testExit <- exitCode
	}()

	select {
	case <-shutdown:
		return types.TestResult{Interrupted: true}, serviceTester.Shutdown()
	case err := <-testError:
		serviceTester.Shutdown() // nolint errcheck
		return types.TestResult{}, err
	case exitCode := <-testExit:
		testResult := types.TestResult{Passed: exitCode == 0}
		result := "failed"
		if testResult.Passed {
			result = "passed"
		}
		logger.Logf("'%s' tests %s", serviceContext.Dir, result)
		return testResult, serviceTester.Shutdown()
	}
}
