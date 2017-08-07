package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runtools"

// RestartContainer starts the docker container of the given service
func RestartContainer(serviceName string, env []string, dockerComposeDir string, logChannel chan string) error {
	return runtools.RunAndLog(dockerComposeDir, env, logChannel, "docker-compose", "restart", serviceName)
}
