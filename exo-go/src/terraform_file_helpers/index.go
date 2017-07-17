package terraformFileHelpers

import (
	"strings"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

// GenerateTerraform generates the main terraform file given application and service configuration
func GenerateTerraform(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig) error {
	fileData := []string{}

	moduleData, err := generateAwsModule(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateServiceModules(serviceConfigs)
	if err != nil {
		return errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	err = WriteTerraformFile(strings.Join(fileData, "\n"))
	if err != nil {
		return errors.Wrap(err, "Failed to write Terraform file")
	}
	return nil
}

func generateAwsModule(appConfig types.AppConfig) (string, error) {
	varsMap := map[string]string{
		"appName": appConfig.Name,
		"region":  "us-west-2", //TODO prompt user for this
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(serviceConfigs map[string]types.ServiceConfig) (string, error) {
	serviceModules := []string{}
	for serviceName, serviceConfig := range serviceConfigs {
		module, err := generateServiceModule(serviceName, serviceConfig)
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceName string, serviceConfig types.ServiceConfig) (string, error) {
	varsMap := map[string]string{
		"serviceRole":    serviceName,
		"startupCommand": serviceConfig.Startup["command"],
		"publicPort":     serviceConfig.Production["publicPort"],
		"cpu":            serviceConfig.Production["cpu"],
	}
	return RenderTemplates("public_service.tf", varsMap)
}
