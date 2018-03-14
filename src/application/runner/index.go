package runner

import (
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/Originate/exosphere/src/docker/composerunner"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/util"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

// Run runs the application with graceful shutdown
func Run(options RunOptions) error {
	err := copySharedDirectories(options.AppContext)
	if err != nil {
		return err
	}
	envVars := map[string]string{
		"COMPOSE_PROJECT_NAME": options.DockerComposeProjectName,
		"APP_PATH":             options.AppContext.Location,
	}
	util.Merge(envVars, buildSecretEnvVars(options.AppContext))
	runOptions := composerunner.RunOptions{
		DockerComposeDir:      options.AppContext.GetDockerComposeDir(),
		DockerComposeFileName: options.DockerComposeFileName,
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
		err = composerunner.Run(runOptions)
		doneChannel <- true
	}()
	<-doneChannel
	_ = composerunner.Shutdown(runOptions)
	if strings.Contains(err.Error(), "exit status") {
		return nil
	}
	return err
}

func copySharedDirectories(appContext *context.AppContext) error {
	for _, serviceSource := range appContext.Config.Services {
		if serviceSource.Location != "" {
			for _, sharedDirectory := range appContext.Config.SharedDirectories {
				dirBase := path.Base(sharedDirectory)
				srcDir := path.Join(appContext.Location, sharedDirectory)
				destDir := path.Join(appContext.Location, serviceSource.Location, dirBase)
				err := osutil.CopyRecursively(srcDir, destDir)
				if err != nil {
					return err
				}
			}
		}
	}
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
