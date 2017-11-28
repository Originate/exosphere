package runner

import (
	"os"
	"os/signal"
	"path"

	"github.com/Originate/exosphere/src/application"
	"github.com/Originate/exosphere/src/docker/composerunner"
)

// Run runs the application with graceful shutdown
func Run(options RunOptions) error {
	err := application.GenerateComposeFiles(options.AppContext)
	if err != nil {
		return err
	}
	runOptions := composerunner.RunOptions{
		AppDir:                   options.AppContext.Location,
		DockerComposeDir:         path.Join(options.AppContext.Location, "docker-compose"),
		DockerComposeFileName:    options.BuildMode.GetDockerComposeFileName(),
		DockerComposeProjectName: options.DockerComposeProjectName,
		Writer: options.Writer,
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
