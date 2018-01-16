package terraform

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/pkg/errors"
)

// GenerateServices generates the contents of the main terraform file given application and service configuration
func GenerateServices(deployConfig deploy.Config) (string, error) {
	fileData := []string{}
	moduleData, err := generateMainServiceModule(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate main service Terraform module")
	}
	fileData = append(fileData, moduleData)
	moduleData, err = generateServiceModules(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)
	return strings.Join(fileData, "\n"), nil
}

func generateMainServiceModule(deployConfig deploy.Config) (string, error) {
	varsMap := map[string]string{
		"appName":          deployConfig.AppContext.Config.Name,
		"lockTable":        deployConfig.AwsConfig.TerraformLockTable,
		"terraformVersion": TerraformVersion,
	}
	return RenderTemplates("services.tf", varsMap)
}

func generateServiceModules(deployConfig deploy.Config) (string, error) {
	serviceModules := []string{}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		serviceConfig := deployConfig.AppContext.ServiceContexts[serviceRole].Config
		module, err := generateServiceModule(serviceRole, deployConfig, serviceConfig, fmt.Sprintf("%s_service.tf", serviceConfig.Type))
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceRole string, deployConfig deploy.Config, serviceConfig types.ServiceConfig, filename string) (string, error) {
	varsMap := map[string]string{
		"serviceRole":         serviceRole,
		"publicPort":          serviceConfig.Production.Port,
		"cpu":                 serviceConfig.Remote.CPU,
		"memory":              serviceConfig.Remote.Memory,
		"healthCheck":         serviceConfig.Production.HealthCheck,
		"terraformCommitHash": TerraformModulesRef,
	}
	return RenderTemplates(filename, varsMap)
}
