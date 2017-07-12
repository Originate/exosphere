package serviceConfigHelpers

import (
	"io/ioutil"
	"log"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/types"
)

// GetServiceConfig reads service.yml of the role service and returns
// the serviceConfig object
func GetServiceConfig(location string) types.ServiceConfig {
	yamlFile, err := ioutil.ReadFile(path.Join(location, "service.yml"))
	if err != nil {
		log.Fatalf("Failed to read service.yml: %s", err)
	}
	var serviceConfig types.ServiceConfig
	err = yaml.Unmarshal(yamlFile, &serviceConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal service.yml for the service at %s: %s", location, err)
	}
	return serviceConfig
}
