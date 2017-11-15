package application

import (
	"io"
	"os"
	"os/signal"
	"path"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types"
)

// Runner runs the overall application
type Runner struct {
	AppContext               types.AppContext
	ServiceConfigs           map[string]types.ServiceConfig
	BuiltDependencies        map[string]config.AppDevelopmentDependency
	DockerComposeDir         string
	DockerComposeProjectName string
	Writer                   io.Writer
	BuildMode                composebuilder.BuildMode
}

// NewRunner is Runner's constructor
func NewRunner(appContext types.AppContext, writer io.Writer, dockerComposeProjectName string, buildMode composebuilder.BuildMode) (*Runner, error) {
	serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
	if err != nil {
		return &Runner{}, err
	}
	allBuiltDependencies := config.GetBuiltDevelopmentDependencies(appContext.Config, serviceConfigs, appContext.Location)
	return &Runner{
		AppContext:               appContext,
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
	dockerCompose, err := composebuilder.GetApplicationDockerCompose(composebuilder.ApplicationOptions{
		AppConfig: r.AppContext.Config,
		AppDir:    r.AppContext.Location,
		BuildMode: r.BuildMode,
	})
	if err != nil {
		return err
	}
	runOptions := composerunner.RunOptions{
		AppDir:                   r.AppContext.Location,
		DockerCompose:            dockerCompose,
		DockerComposeDir:         r.DockerComposeDir,
		DockerComposeFileName:    r.BuildMode.GetDockerComposeFileName(),
		DockerComposeProjectName: r.DockerComposeProjectName,
		Writer: r.Writer,
	}
	doneChannel := make(chan bool, 1)
	go func() {
		sigIntChannel := make(chan os.Signal, 1)
		signal.Notify(sigIntChannel, os.Interrupt)
		<-sigIntChannel
		signal.Stop(sigIntChannel)
		doneChannel <- true
	}()
	go func() {
		_ = composerunner.Run(runOptions)
		doneChannel <- true
	}()
	<-doneChannel
	_ = composerunner.Shutdown(runOptions)
	return nil
}
