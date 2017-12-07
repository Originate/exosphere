package composerunner

import (
	"fmt"

	"github.com/Originate/exosphere/src/docker/compose"
)

// Run runs docker images based on the given options
func Run(options RunOptions) error {
	err := compose.RunImages(compose.CommandOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   buildEnvSlice(options.EnvironmentVariables),
		AbortOnExit:           options.AbortOnExit,
		Build:                 true,
	})
	return err
}

// RunService runs a service based on the given options
func RunService(options RunOptions, serviceName string) error {
	err := compose.RunImages(compose.CommandOptions{
		DockerComposeDir:      options.DockerComposeDir,
		DockerComposeFileName: options.DockerComposeFileName,
		Writer:                options.Writer,
		Env:                   buildEnvSlice(options.EnvironmentVariables),
		AbortOnExit:           options.AbortOnExit,
		Build:                 true,
		ImageNames:            []string{serviceName},
	})
	return err
}

func buildEnvSlice(envMap map[string]string) []string {
	env := []string{}
	for k, v := range envMap {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}
