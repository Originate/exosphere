package application

import (
	"io"
	"os"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
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
			util.PrintSectionHeaderf(writer, "%s has no tests, skipping\n", serviceContext.Dir)
		} else {
			testResult, err := TestService(serviceContext, writer, mode, shutdown)
			if err != nil {
				util.PrintSectionHeaderf(writer, "error running '%s' tests:", err)
			}
			if testResult.Interrupted {
				return testResult, nil
			}
			if !testResult.Passed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		util.PrintSectionHeader(writer, "All tests passed\n\n")
		return types.TestResult{
			Passed: true,
		}, nil
	}
	util.PrintSectionHeaderf(writer, "%d tests failed\n\n", numFailed)
	return types.TestResult{}, nil
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext types.ServiceContext, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	util.PrintSectionHeaderf(writer, "Testing service '%s'\n", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, writer, mode)
	if err != nil {
		return types.TestResult{}, err
	}

	testExit := make(chan int)
	testError := make(chan error)
	go func() {
		exitCode, err := serviceTester.Run()
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
		return types.TestResult{Passed: exitCode == 0}, serviceTester.Shutdown()
	}
}
