package composebuilder

import (
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	yaml "gopkg.in/yaml.v2"
)

// WriteYML writes a docker-compose.yml file
func WriteYML(dir string, dockerConfigs types.DockerConfigs) error {
	bytes, err := yaml.Marshal(types.DockerCompose{
		Version:  "3",
		Services: dockerConfigs,
	})
	if err != nil {
		return err
	}
	if err := util.CreateEmptyDirectory(dir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dir, "docker-compose.yml"), bytes, 0777)
}
