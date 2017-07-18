package terraformFileHelpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/app_dependency_helpers"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

// GenerateTerraform generates the main terraform file given application and service configuration
func GenerateTerraform(appConfig types.AppConfig, serviceConfigs map[string]types.ServiceConfig, appDir string) error {
	fileData := []string{}

	moduleData, err := generateAwsModule(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	serviceProtectionLevels := appConfigHelpers.GetServiceProtectionLevels(appConfig)
	moduleData, err = generateServiceModules(serviceConfigs, serviceProtectionLevels)
	if err != nil {
		return errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(appConfig)
	if err != nil {
		return errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	err = WriteTerraformFile(strings.Join(fileData, "\n"), appDir)
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

func generateServiceModules(serviceConfigs map[string]types.ServiceConfig, serviceProtectionLevels map[string]string) (string, error) {
	serviceModules := []string{}
	for serviceName, serviceConfig := range serviceConfigs {
		var module string
		var err error
		switch serviceProtectionLevels[serviceName] {
		case "public":
			module, err = generateServiceModule(serviceName, serviceConfig, "public_service.tf")
		case "private":
			module, err = generateServiceModule(serviceName, serviceConfig, "private_service.tf")
		}
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceName string, serviceConfig types.ServiceConfig, filename string) (string, error) {
	command, err := json.Marshal(strings.Split(serviceConfig.Startup["command"], " "))
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal service startup command")
	}
	varsMap := map[string]string{
		"serviceRole":    serviceName,
		"startupCommand": string(command),
		"publicPort":     serviceConfig.Production["public-port"],
		"cpu":            serviceConfig.Production["cpu"],
		"memory":         serviceConfig.Production["memory"],
		"url":            serviceConfig.Production["url"],
		"healthCheck":    serviceConfig.Production["health-check"],
		//"envVars": TODO: determine how we define env vars and then implement
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(appConfig types.AppConfig) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range appConfig.Dependencies {
		deploymentConfig := appDependencyHelper.Build(dependency, appConfig).GetDeploymentConfig()
		module, err := RenderTemplates(fmt.Sprintf("%s.tf", dependency.Name), deploymentConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}
