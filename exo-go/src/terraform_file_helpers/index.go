package terraformFileHelpers

import (
	"strings"

	"github.com/Originate/exosphere/exo-go/src/types"
)

// GenerateTerraform generates the main terraform file given application and service configuration
func GenerateTerraform(appConfig types.AppConfig) {
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
