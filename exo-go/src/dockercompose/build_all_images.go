package dockercompose

import "github.com/Originate/exosphere/exo-go/src/util"

// BuildAllImages builds all the docker images defined in docker-compose.yml
func BuildAllImages(dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "build")
}
