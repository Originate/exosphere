package terraformFileHelpers

import (
	"strings"

	"github.com/Originate/exosphere/exo-go/src/types"
)

func GenerateTerraform(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig) {
	fileData := []string{}
	fileData = append(fileData, generateAwsModule(appConfig))
	WriteTerraformFile(strings.Join(fileData, "\n"))
}

func generateAwsModule(appConfig types.AppConfig) string {
	varsMap := map[string]string{
		"appName": appConfig.Name,
		"region":  "us-west-2", //TODO prompt user for this
	}
	return RenderTemplates("aws.tf", varsMap)
}
