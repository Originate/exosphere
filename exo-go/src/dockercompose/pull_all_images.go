package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runtools"

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeDir string, logChannel chan string) error {
	return runtools.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "pull")
}
