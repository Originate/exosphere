package appConfigHelpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/src/user_input_helpers"
	"gopkg.in/yaml.v2"
)

// GetAppConfig reads application.yml and returns the appConfig object
func GetAppConfig() types.AppConfig {
	yamlFile, err := ioutil.ReadFile("application.yml")
	if err != nil {
		log.Fatalf("Failed to read application.yml: %s", err)
	}
	var appConfig types.AppConfig
	err = yaml.Unmarshal(yamlFile, &appConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal application.yml: %s", err)
	}
	return appConfig
}

// UpdateAppConfig adds serviceRole to the appConfig object and updates
// application.yml
func UpdateAppConfig(reader *bufio.Reader, serviceRole string, appConfig types.AppConfig) {
	switch userInputHelpers.Choose(reader, "Protection Level:", []string{"public", "private"}) {
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
		log.Fatalf("Failed to marshal application.yml: %s", err)
	}
	err = ioutil.WriteFile(path.Join("application.yml"), bytes, 0777)
	if err != nil {
		log.Fatalf("Failed to write application.yml: %s", err)
	}
}
