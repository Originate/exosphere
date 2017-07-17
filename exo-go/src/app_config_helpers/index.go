package appConfigHelpers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
	"github.com/segmentio/go-prompt"
	"gopkg.in/yaml.v2"
)

// GetAppConfig reads application.yml and returns the appConfig object
func GetAppConfig() (result types.AppConfig, err error) {
	yamlFile, err := ioutil.ReadFile("application.yml")
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return result, errors.Wrap(err, "Failed to unmarshal application.yml")
	}
	return result, nil
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(serviceRole string, appConfig types.AppConfig) error {
	protectionLevels := []string{"public", "private"}
	switch protectionLevels[prompt.Choose("Protection Level:", protectionLevels)] {
	case "public":
		if appConfig.Services.Public == nil {
			appConfig.Services.Public = make(map[string]types.ServiceConfig)
		}
		appConfig.Services.Public[serviceRole] = types.ServiceConfig{Location: fmt.Sprintf("./%s", serviceRole)}
	case "private":
		if appConfig.Services.Private == nil {
			appConfig.Services.Private = make(map[string]types.ServiceConfig)
		}
		appConfig.Services.Private[serviceRole] = types.ServiceConfig{Location: fmt.Sprintf("./%s", serviceRole)}
	}
	bytes, err := yaml.Marshal(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal application.yml")
	}
	return ioutil.WriteFile(path.Join("application.yml"), bytes, 0777)
}
