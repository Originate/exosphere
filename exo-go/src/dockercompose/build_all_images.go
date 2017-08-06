package dockercompose

import "github.com/Originate/exosphere/exo-go/src/runtools"

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeDir string, logChannel chan string) error {
	return runtools.AndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "build")
}
