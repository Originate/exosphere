package tester

import (
	"io"
	"os"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
	"github.com/fatih/color"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext *context.AppContext, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	failedTests := []string{}
	locations := []string{}
	testRunner, err := NewTestRunner(appContext, writer, mode)
	if err != nil {
		return types.TestResult{}, err
	}
	for _, serviceRole := range appContext.Config.GetSortedServiceRoles() {
		serviceContext := appContext.ServiceContexts[serviceRole]
		serviceLocation := serviceContext.Source.Location
		if serviceLocation == "" || util.DoesStringArrayContain(locations, serviceLocation) {
			continue
		}
		locations = append(locations, serviceLocation)
		if serviceContext.Config.Development.Scripts["test"] == "" {
			util.PrintSectionHeaderf(writer, "%s has no tests, skipping\n", serviceContext.ID())
		} else {
			var testResult types.TestResult
			testResult, err = runServiceTest(testRunner, serviceContext, writer, shutdown)
			if err != nil {
				util.PrintSectionHeaderf(writer, "error running '%s' tests:", err)
			}
			if testResult.Interrupted {
				return testResult, nil
			}
			if !testResult.Passed {
				failedTests = append(failedTests, serviceContext.ID())
			}
		}
	}
	err = printResults(failedTests, writer)
	return types.TestResult{
		Passed: len(failedTests) == 0,
	}, err
}

func printResults(failedTests []string, writer io.Writer) error {
	if len(failedTests) == 0 {
		green := color.New(color.FgGreen)
		_, err := green.Fprint(writer, "All tests passed\n\n")
		if err != nil {
			return err
		}
	}
	red := color.New(color.FgRed)
	_, err := red.Fprint(writer, "The following tests failed:\n")
	if err != nil {
		return err
	}
	for _, failedTest := range failedTests {
		_, err = red.Fprintln(writer, failedTest)
		if err != nil {
			return err
		}
	}
	return nil
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext *context.ServiceContext, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	testRunner, err := NewTestRunner(serviceContext.AppContext, writer, mode)
	if err != nil {
		return types.TestResult{}, err
	}
	return runServiceTest(testRunner, serviceContext, writer, shutdown)
}

func runServiceTest(testRunner *TestRunner, serviceContext *context.ServiceContext, writer io.Writer, shutdown chan os.Signal) (types.TestResult, error) {
	util.PrintSectionHeaderf(writer, "Testing service '%s'\n", serviceContext.ID())

	testExit := make(chan int)
	testError := make(chan error)
	go func() {
		exitCode, err := testRunner.RunTest(serviceContext.Role)
		if err != nil {
			testError <- err
			return
		}
		testExit <- exitCode
	}()

	select {
	case <-shutdown:
		return types.TestResult{Interrupted: true}, testRunner.Shutdown()
	case err := <-testError:
		testRunner.Shutdown() // nolint errcheck
		return types.TestResult{}, err
	case exitCode := <-testExit:
		return types.TestResult{Passed: exitCode == 0}, testRunner.Shutdown()
	}
}
