package dockercompose

import "github.com/Originate/exosphere/exo-go/src/util"

// PullAllImages pulls all the docker images defined in docker-compose.yml
func PullAllImages(dockerComposeDir string, logChannel chan string) error {
	return util.RunAndLog(dockerComposeDir, []string{}, logChannel, "docker-compose", "pull")
}
