package dockercompose

import "github.com/Originate/exosphere/exo-go/src/util"

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(serviceName string, env []string, dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, env, logChannel, "docker-compose", "create", "--build", serviceName)
}
