package tester

import (
	"io"
	"strings"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

// TestRunner runs the tests for the given service
type TestRunner struct {
	AppContext *context.AppContext
	BuildMode  types.BuildMode
	RunOptions composerunner.RunOptions
	Writer     io.Writer
}

// NewTestRunner is TestRunner's constructor
func NewTestRunner(appContext *context.AppContext, writer io.Writer) (*TestRunner, error) {
	tester := &TestRunner{
		AppContext: appContext,
		Writer:     writer,
	}
	var err error
	tester.RunOptions, err = tester.getRunOptions()
	if err != nil {
		return nil, err
	}
	return tester, err
}

// RunTest runs the tests for the service and return true if the tests passed and an error if any
func (s *TestRunner) RunTest(serviceRole string) (int, error) {
	err := composerunner.RunService(s.RunOptions, serviceRole)
	if err != nil {
		if strings.Contains(err.Error(), "exit status") {
			return 1, nil
		}
		return 1, err
	}
	return 0, nil
}

// Shutdown shuts down the tests
func (s *TestRunner) Shutdown() error {
	return composerunner.Shutdown(s.RunOptions)
}

func (s *TestRunner) getRunOptions() (composerunner.RunOptions, error) {
	dockerComposeProjectName := composebuilder.GetTestDockerComposeProjectName(s.AppContext.Config.Name)
	return composerunner.RunOptions{
		DockerComposeDir:      s.AppContext.GetDockerComposeDir(),
		DockerComposeFileName: types.LocalTestComposeFileName,
		Writer:                s.Writer,
		AbortOnExit:           true,
		EnvironmentVariables: map[string]string{
			"COMPOSE_PROJECT_NAME": dockerComposeProjectName,
			"APP_PATH":             s.AppContext.Location,
		},
	}, nil
}
