package appSetup

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/test_helpers"
	"github.com/pkg/errors"
)

func joinDockerConfigMaps(maps ...map[string]types.DockerConfig) map[string]types.DockerConfig {
	result := map[string]types.DockerConfig{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

// CheckoutApp copies the example app appName to cwd
func CheckoutApp(cwd, appName string) error {
	src := path.Join("..", "..", "..", "exosphere-shared", "example-apps", appName)
	dest := path.Join(cwd, appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return testHelpers.CopyDir(src, dest)
}

// GetDockerCompose reads docker-compose.yml at the given location and
// returns the dockerCompose object
func GetDockerCompose(dockerComposeLocation string) (result types.DockerCompose, err error) {
	yamlFile, err := ioutil.ReadFile(dockerComposeLocation)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "Failed to unmarshal docker-compose.yml")
	}
	return result, nil
}
