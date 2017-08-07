package dockercompose

import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// GetDockerCompose reads docker-compose.yml at the given path and
// returns the dockerCompose object
func GetDockerCompose(dockerComposePath string) (result DockerCompose, err error) {
	yamlFile, err := ioutil.ReadFile(dockerComposePath)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "Failed to unmarshal docker-compose.yml")
	}
	return result, nil
}
