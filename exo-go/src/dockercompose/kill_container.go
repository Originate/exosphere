package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runtools"

// KillContainer kills the docker container of the given service
func KillContainer(serviceName, dockerComposeDir string, logChannel chan string) error {
	return runtools.AndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "kill", serviceName)
}
