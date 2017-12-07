package runner

import (
	"os"
	"os/signal"
	"path"

	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
)

// Run runs the application with graceful shutdown
func Run(options RunOptions) error {
	envVars := map[string]string{
		"COMPOSE_PROJECT_NAME": options.DockerComposeProjectName,
		"APP_PATH":             options.AppContext.Location,
	}
	util.Merge(envVars, buildSecretEnvVars(options.AppContext))
	runOptions := composerunner.RunOptions{
		DockerComposeDir:      path.Join(options.AppContext.Location, "docker-compose"),
		DockerComposeFileName: options.BuildMode.GetDockerComposeFileName(),
		Writer:                options.Writer,
		EnvironmentVariables:  envVars,
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

func buildSecretEnvVars(appContext *context.AppContext) map[string]string {
	secrets := map[string]string{}
	for _, serviceContext := range appContext.ServiceContexts {
		for _, secretName := range serviceContext.Config.Local.Secrets {
			secrets[secretName] = os.Getenv(secretName)
		}
	}
	return secrets
}
