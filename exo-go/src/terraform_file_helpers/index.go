package terraformFileHelpers

import (
	"strings"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

// GenerateTerraform generates the main terraform file given application and service configuration
func GenerateTerraform(appConfig types.AppConfig) error {
	fileData := []string{}

	moduleData, err := generateAwsModule(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to generate Terraform")
	}
	fileData = append(fileData, moduleData)

	WriteTerraformFile(strings.Join(fileData, "\n"))
	return nil
}

func generateAwsModule(appConfig types.AppConfig) (string, error) {
	varsMap := map[string]string{
		"appName": appConfig.Name,
		"region":  "us-west-2", //TODO prompt user for this
	}
	return RenderTemplates("aws.tf", varsMap)
}
