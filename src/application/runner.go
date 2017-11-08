package application

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
)

// Runner runs the overall application
type Runner struct {
	AppContext               types.AppContext
	HomeDir                  string
	ServiceConfigs           map[string]types.ServiceConfig
	BuiltDependencies        map[string]config.AppDevelopmentDependency
	DockerComposeDir         string
	DockerComposeProjectName string
	Writer                   io.Writer
	BuildMode                composebuilder.BuildMode
}

// NewRunner is Runner's constructor
func NewRunner(appContext types.AppContext, writer io.Writer, homeDir, dockerComposeProjectName string, buildMode composebuilder.BuildMode) (*Runner, error) {
	serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
	if err != nil {
		return &Runner{}, err
	}
	allBuiltDependencies := config.GetBuiltDevelopmentDependencies(appContext.Config, serviceConfigs, appContext.Location, homeDir)
	return &Runner{
		AppContext:               appContext,
		HomeDir:                  homeDir,
		ServiceConfigs:           serviceConfigs,
		BuiltDependencies:        allBuiltDependencies,
		DockerComposeDir:         path.Join(appContext.Location, "docker-compose"),
		DockerComposeProjectName: dockerComposeProjectName,
		Writer:    writer,
		BuildMode: buildMode,
	}, nil
}

// Run runs the application with graceful shutdown
func (r *Runner) Run() error {
	dockerConfigs, err := composebuilder.GetApplicationDockerConfigs(composebuilder.ApplicationOptions{
		AppConfig: r.AppContext.Config,
		AppDir:    r.AppContext.Location,
		BuildMode: r.BuildMode,
		HomeDir:   r.HomeDir,
	})
	if err != nil {
		return err
	}
	runOptions := composerunner.RunOptions{
		DockerConfigs:            dockerConfigs,
		DockerComposeDir:         r.DockerComposeDir,
		DockerComposeFileName:    r.BuildMode.GetDockerComposeFileName(),
		DockerComposeProjectName: r.DockerComposeProjectName,
		Writer: r.Writer,
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		shutdownChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownChannel, os.Interrupt)
		<-shutdownChannel
		signal.Stop(shutdownChannel)
		if shutdownErr := composerunner.Shutdown(runOptions); shutdownErr != nil {
			fmt.Println("Error shutting down")
		}
		wg.Done()
	}()
	err = composerunner.Run(runOptions)
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		_ = composerunner.Shutdown(runOptions)
		return err
	}
	wg.Wait()
	return nil
}
