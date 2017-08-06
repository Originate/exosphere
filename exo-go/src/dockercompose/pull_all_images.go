package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runplus"

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeDir string, logChannel chan string) error {
	return runplus.AndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "pull")
}
