package composebuilder

import (
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/util"
	yaml "gopkg.in/yaml.v2"
)

// WriteYML writes a docker-compose.yml file
func WriteYML(dir, filename string, file *types.DockerCompose) error {
	bytes, err := yaml.Marshal(file)
	if err != nil {
		return err
	}
	if err := util.MakeDirectory(dir); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dir, filename), bytes, 0777)
}
