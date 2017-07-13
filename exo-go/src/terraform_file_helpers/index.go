package terraformFileHelpers

import (
	"github.com/Originate/exosphere/exo-go/src/types"
)

func GenerateTerraform(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig) {
	generateAwsModule(appConfig)
}

func generateAwsModule(appConfig types.AppConfig) {
	varsMap := map[string]string{
		"appName": appConfig.Name,
		"region":  "us-west-2", //TODO prompt user for this
	}
	RenderTemplates("aws.tf", varsMap)
}
