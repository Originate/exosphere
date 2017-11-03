package tester

import (
	"io"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
)

// TestRunner runs the tests for the given service
type TestRunner struct {
	AppConfig                types.AppConfig
	AppDir                   string
	BuildMode                composebuilder.BuildMode
	DockerComposeProjectName string
	HomeDir                  string
	RunOptions               composerunner.RunOptions
	Writer                   io.Writer
}

// NewTestRunner is TestRunner's constructor
func NewTestRunner(appContext types.AppContext, writer io.Writer, mode composebuilder.BuildMode) (*TestRunner, error) {
	appConfig := appContext.Config
	appDir := appContext.Location
	homeDir, err := util.GetHomeDirectory()
	if err != nil {
		return nil, err
	}
	dockerComposeProjectName := composebuilder.GetTestDockerComposeProjectName(appConfig.Name)

	tester := &TestRunner{
		AppConfig:                appConfig,
		AppDir:                   appDir,
		BuildMode:                mode,
		DockerComposeProjectName: dockerComposeProjectName,
		HomeDir:                  homeDir,
		Writer:                   writer,
	}
	tester.RunOptions, err = tester.getRunOptions()
	if err != nil {
		return nil, err
	}
	return tester, err
}

// RunTest runs the tests for the service and return true if the tests passed and an error if any
func (s *TestRunner) RunTest(testRole string) (int, error) {
	err := composerunner.RunService(s.RunOptions, testRole)
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
	dockerComposeConfigs, err := s.getDockerComposeConfigs()
	if err != nil {
		return composerunner.RunOptions{}, err
	}
	return composerunner.RunOptions{
		DockerConfigs:            dockerComposeConfigs,
		DockerComposeDir:         path.Join(s.AppDir, "docker-compose"),
		DockerComposeFileName:    s.BuildMode.GetDockerComposeFileName(),
		DockerComposeProjectName: s.DockerComposeProjectName,
		Writer:      s.Writer,
		AbortOnExit: true,
	}, nil
}

func (s *TestRunner) getDockerComposeConfigs() (types.DockerConfigs, error) {
	return composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: s.AppConfig,
		AppDir:    s.AppDir,
		BuildMode: s.BuildMode,
		HomeDir:   s.HomeDir,
	})
}
