package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runplus"

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeDir string, logChannel chan string) error {
	return runplus.AndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "build")
}
