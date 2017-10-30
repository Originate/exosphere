package application

import (
	"io"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestApp runs the tests for the entire application and return true if the tests passed
// and an error if any
func TestApp(appContext types.AppContext, writer io.Writer, mode composebuilder.BuildMode) (bool, error) {
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
			util.PrintSectionHeaderf(writer, "%s has no tests, skipping\n", serviceContext.Dir)
		} else {
			testPassed, err := TestService(serviceContext, writer, mode)
			if err != nil {
				util.PrintSectionHeaderf(writer, "error running '%s' tests:", err)
			}
			if !testPassed {
				numFailed++
			}
		}
	}
	if numFailed == 0 {
		util.PrintSectionHeader(writer, "All tests passed\n\n")
		return true, nil
	}
	util.PrintSectionHeaderf(writer, "%d tests failed\n\n", numFailed)
	return false, nil
}

// TestService runs the tests for the service and return true if the tests passed
// and an error if any
func TestService(serviceContext types.ServiceContext, writer io.Writer, mode composebuilder.BuildMode) (bool, error) {
	util.PrintSectionHeaderf(writer, "Testing service '%s'\n", serviceContext.Dir)
	serviceTester, err := NewServiceTester(serviceContext, writer, mode)
	if err != nil {
		return false, err
	}

	exitCode, err := serviceTester.Run()
	if err != nil {
		return false, err
	}
	return exitCode == 0, serviceTester.Shutdown()
}
