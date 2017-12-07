package composerunner

import (
	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/pkg/errors"
)

// Shutdown kills the docker images based on the given options
func Shutdown(options RunOptions) error {
	err := compose.KillContainers(compose.CommandOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   buildEnvSlice(options.EnvironmentVariables),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to shutdown the app")
	}
	return nil
}
