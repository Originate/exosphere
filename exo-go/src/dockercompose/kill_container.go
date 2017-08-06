package dockercompose

import "github.com/Originate/exosphere/exo-go/src/run"

// KillContainer kills the docker container of the given service
func KillContainer(serviceName, dockerComposeDir string, logChannel chan string) error {
	return run.AndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "kill", serviceName)
}
