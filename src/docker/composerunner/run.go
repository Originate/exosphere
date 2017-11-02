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
	err = buildAndPullImages(options)
	if err != nil {
		return err
	}
	err = compose.RunImages(compose.ImagesOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
		AbortOnExit:           options.AbortOnExit,
	})
	return err
}

func buildAndPullImages(options RunOptions) error {
	opts := compose.BaseOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
	}
	err := compose.PullAllImages(opts)
	if err != nil {
		return err
	}
	return compose.BuildAllImages(opts)
}
