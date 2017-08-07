package terraform

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/config"
	"github.com/pkg/errors"
)

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(config config.TerraformConfig) error {
	fileData, err := Generate(config)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, config.TerraformDir)
	return err
}

// Generate generates the contents of the main terraform file given application and service configuration
func Generate(config config.TerraformConfig) (string, error) {
	fileData := []string{}

	moduleData, err := generateAwsModule(config)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	serviceProtectionLevels := config.AppConfig.GetServiceProtectionLevels()
	moduleData, err = generateServiceModules(config.ServiceConfigs, serviceProtectionLevels)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(config)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	return strings.Join(fileData, "\n"), nil
}

func generateAwsModule(config config.TerraformConfig) (string, error) {
	varsMap := map[string]string{
		"appName":      config.AppConfig.Name,
		"remoteBucket": config.RemoteBucket,
		"lockTable":    config.LockTable,
		"region":       config.Region,
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(serviceConfigs map[string]config.ServiceConfig, serviceProtectionLevels map[string]string) (string, error) {
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

func generateServiceModule(serviceName string, serviceConfig config.ServiceConfig, filename string) (string, error) {
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

func generateDependencyModules(terraformConfig config.TerraformConfig) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range terraformConfig.AppConfig.Dependencies {
		deploymentConfig := config.NewAppDependency(dependency, terraformConfig.AppConfig, terraformConfig.AppDir, terraformConfig.HomeDir).GetDeploymentConfig()
		module, err := RenderTemplates(fmt.Sprintf("%s.tf", dependency.Name), deploymentConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}
