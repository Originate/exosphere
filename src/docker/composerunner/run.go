package composerunner

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/Originate/exosphere/src/docker/composebuilder"
)

// Run runs docker images based on the given options
func Run(options RunOptions) error {
	err := composebuilder.WriteYML(options.DockerComposeDir, options.DockerComposeFileName, options.DockerConfigs)
	if err != nil {
		return err
	}
	err = compose.RunImages(compose.CommandOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env: []string{
			fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName),
			fmt.Sprintf("APP_PATH=%s", options.AppDir),
		},
		AbortOnExit: options.AbortOnExit,
		Build:       true,
	})
	return err
}

// RunService runs a service based on the given options
func RunService(options RunOptions, serviceName string) error {
	err := composebuilder.WriteYML(options.DockerComposeDir, options.DockerComposeFileName, options.DockerConfigs)
	if err != nil {
		return err
	}
	err = compose.RunImages(compose.CommandOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env: []string{
			fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName),
			fmt.Sprintf("APP_PATH=%s", options.AppDir),
		},
		AbortOnExit: options.AbortOnExit,
		Build:       true,
		ImageNames:  []string{serviceName},
	})
	return err
}
