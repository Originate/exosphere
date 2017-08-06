package config

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func getInternalServiceConfig(appDir, serviceDirName string) (ServiceConfig, error) {
	var serviceConfig ServiceConfig
	yamlFile, err := ioutil.ReadFile(path.Join(appDir, serviceDirName, "service.yml"))
	if err != nil {
		return serviceConfig, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return serviceConfig, errors.Wrap(err, fmt.Sprintf("Failed to unmarshal service.yml for the internal service '%s'", serviceDirName))
	}
	return serviceConfig, nil
}
