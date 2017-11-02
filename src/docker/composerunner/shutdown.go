package composerunner

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
	"github.com/pkg/errors"
)

// Shutdown kills the docker images based on the given options
func Shutdown(options RunOptions) error {
	err := composebuilder.WriteYML(options.DockerComposeDir, options.DockerComposeFileName, options.DockerConfigs)
	if err != nil {
		return err
	}
	return killImages(options)
}

func killImages(options RunOptions) error {
	err := compose.KillAllContainers(compose.BaseOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
	})
	if err != nil {
		return errors.Wrap(err, "Failed to shutdown the app")
	}
	return nil
}
