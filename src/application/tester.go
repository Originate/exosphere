package application

import (
	"os"
	"os/signal"
	"sync"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed, if they were interrupted,
// and an error if any
func TestApp(appContext types.AppContext, logger *util.Logger, mode composebuilder.BuildMode) (bool, bool, error) {
	logger.Logf("Testing application %s", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return false, false, err
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
			testPassed, testInterrupted, err := TestService(serviceContext, logger, mode)
			if testInterrupted {
				return false, true, err
			}
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
		return true, false, nil
	}
	logger.Logf("%d tests failed", numFailed)
	return false, false, nil
}

// TestService runs the tests for the service and return true if the tests passed, if they were interrupted,
// and an error if any
func TestService(serviceContext types.ServiceContext, logger *util.Logger, mode composebuilder.BuildMode) (bool, bool, error) {
	logger.Logf("Testing service '%s'", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, logger, mode)
	if err != nil {
		return false, false, err
	}

	var isInterrupted bool
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		signal.Stop(c)
		err = serviceTester.Shutdown()
		isInterrupted = true
		wg.Done()
	}()
	if isInterrupted {
		return false, true, err
	}

	exitCode, err := serviceTester.Run()
	if err != nil {
		return false, false, err
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
	return testPassed, false, serviceTester.Shutdown()
}
