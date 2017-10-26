package application

import (
	"fmt"
	"io"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, writer io.Writer, mode composebuilder.BuildMode) (bool, error) {
	fmt.Fprintf(writer, "Testing application %s\n", appContext.Config.Name)
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return false, err
	}

	numFailed := 0
	locations := []string{}
	for _, serviceContext := range serviceContexts {
		if util.DoesStringArrayContain(locations, serviceContext.Location) {
			continue
		}
		locations = append(locations, serviceContext.Location)
		if serviceContext.Config.Development.Scripts["test"] == "" {
			fmt.Fprintf(writer, "%s has no tests, skipping\n", serviceContext.Dir)
		} else {
			testPassed, err := TestService(serviceContext, writer, mode)
			if err != nil {
				fmt.Fprintf(writer, "error running '%s' tests:", err)

			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		fmt.Fprintln(writer, "All tests passed")
		return true, nil
	}
	fmt.Fprintf(writer, "%d tests failed\n", numFailed)
	return false, nil
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext types.ServiceContext, writer io.Writer, mode composebuilder.BuildMode) (bool, error) {
	fmt.Fprintf(writer, "Testing service '%s'\n", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, writer, mode)
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
	fmt.Fprintf(writer, "'%s' tests %s\n", serviceContext.Dir, result)
	return testPassed, serviceTester.Shutdown()
}
