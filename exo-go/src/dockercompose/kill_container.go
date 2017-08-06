package dockercompose

import "github.com/Originate/exosphere/exo-go/src/util"

// KillContainer kills the docker container of the given service
func KillContainer(serviceName, dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "kill", serviceName)
}
