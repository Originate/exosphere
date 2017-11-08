package tester

import (
	"io"
	"os"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	"github.com/fatih/color"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return types.TestResult{}, err
	}

	failedTests := []string{}
	locations := []string{}
	testRunner, err := NewTestRunner(appContext, writer, mode)
	if err != nil {
		return types.TestResult{}, err
	}
	for _, serviceContext := range serviceContexts {
		if util.DoesStringArrayContain(locations, serviceContext.Location) {
			continue
		}
		locations = append(locations, serviceContext.Location)
		if serviceContext.Config.Development.Scripts["test"] == "" {
			util.PrintSectionHeaderf(writer, "%s has no tests, skipping\n", serviceContext.Dir)
		} else {
			testRole := path.Base(serviceContext.Location)
			testResult, err := runServiceTest(testRunner, testRole, writer, shutdown)
			if err != nil {
				util.PrintSectionHeaderf(writer, "error running '%s' tests:", err)
			}
			if testResult.Interrupted {
				return testResult, nil
			}
			if !testResult.Passed {
				failedTests = append(failedTests, testRole)
			}
		}
	}
	err = printResults(failedTests, writer)
	return types.TestResult{
		Passed: len(failedTests) == 0,
	}, err
}

func printResults(failedTests []string, writer io.Writer) error {
	var err error
	if len(failedTests) == 0 {
		green := color.New(color.FgGreen)
		_, err = green.Fprint(writer, "All tests passed\n\n")
		if err != nil {
			return err
		}
	}
	red := color.New(color.FgRed)
	_, err = red.Fprint(writer, "The following tests failed:\n")
	for _, failedTest := range failedTests {
		_, err = red.Fprintln(writer, failedTest)
		if err != nil {
			return err
		}
	}
	return err
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(context types.Context, writer io.Writer, mode composebuilder.BuildMode, shutdown chan os.Signal) (types.TestResult, error) {
	testRunner, err := NewTestRunner(context.AppContext, writer, mode)
	if err != nil {
		return types.TestResult{}, err
	}
	testRole := path.Base(context.ServiceContext.Location)
	return runServiceTest(testRunner, testRole, writer, shutdown)
}

func runServiceTest(testRunner *TestRunner, testRole string, writer io.Writer, shutdown chan os.Signal) (types.TestResult, error) {
	util.PrintSectionHeaderf(writer, "Testing service '%s'\n", testRole)

	testExit := make(chan int)
	testError := make(chan error)
	go func() {
		exitCode, err := testRunner.RunTest(testRole)
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
