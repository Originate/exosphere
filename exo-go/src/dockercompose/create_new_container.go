package dockercompose

import "github.com/Originate/exosphere/exo-go/src/run"

// CreateNewContainer creates a new docker container for the given service
func CreateNewContainer(serviceName string, env []string, dockerComposeDir string, logChannel chan string) error {
	return run.AndLog(dockerComposeDir, env, logChannel, "docker-compose", "create", "--build", serviceName)
}
