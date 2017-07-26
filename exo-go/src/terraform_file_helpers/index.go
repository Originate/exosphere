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

// GenerateTerraformFile generates the main terraform file given application and service configuration
func GenerateTerraformFile(config types.TerraformConfig) error {

	fileData, err := GenerateTerraform(config)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, config.TerraformDir)
	return err
}

// GenerateTerraform generates the contents of the main terraform file given application and service configuration
func GenerateTerraform(config types.TerraformConfig) (string, error) {
	fileData := []string{}

	moduleData, err := generateAwsModule(config.AppConfig, config.RemoteBucket, config.LockTable, config.Region)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	serviceProtectionLevels := appConfigHelpers.GetServiceProtectionLevels(config.AppConfig)
	moduleData, err = generateServiceModules(config.ServiceConfigs, serviceProtectionLevels)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(config.AppConfig, config.AppDir, config.HomeDir)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	return strings.Join(fileData, "\n"), nil
}

func generateAwsModule(appConfig types.AppConfig, remoteBucket, lockTable, region string) (string, error) {
	varsMap := map[string]string{
		"appName":      appConfig.Name,
		"remoteBucket": remoteBucket,
		"lockTable":    lockTable,
		"region":       region,
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
		//"dockerImage": TODO: implement after ecr functionality is in place
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(appConfig types.AppConfig, appDir, homeDir string) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range appConfig.Dependencies {
		deploymentConfig := appDependencyHelpers.Build(dependency, appConfig, appDir, homeDir).GetDeploymentConfig()
		module, err := RenderTemplates(fmt.Sprintf("%s.tf", dependency.Name), deploymentConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}
