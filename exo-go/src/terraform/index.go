package terraform

import (
	"encoding/json"
	"fmt"
	"strings"

	appDependency "github.com/Originate/exosphere/exo-go/src/config"
	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/pkg/errors"
)

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(config types.DeployConfig, imagesMap map[string]string) error {
	fileData, err := Generate(config, imagesMap)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, config.TerraformDir)
	return err
}

// Generate generates the contents of the main terraform file given application and service configuration
func Generate(config types.DeployConfig, imagesMap map[string]string) (string, error) {
	fileData := []string{}

	moduleData, err := generateAwsModule(config)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	serviceProtectionLevels := config.AppConfig.GetServiceProtectionLevels()
	moduleData, err = generateServiceModules(config.ServiceConfigs, serviceProtectionLevels, imagesMap)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(config, imagesMap)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	return strings.Join(fileData, "\n"), nil
}

func generateAwsModule(config types.DeployConfig) (string, error) {
	varsMap := map[string]string{
		"appName":     config.AppConfig.Name,
		"stateBucket": config.AwsConfig.TerraformStateBucket,
		"lockTable":   config.AwsConfig.TerraformLockTable,
		"region":      config.AwsConfig.Region,
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(serviceConfigs map[string]types.ServiceConfig, serviceProtectionLevels, imagesMap map[string]string) (string, error) {
	serviceModules := []string{}
	for serviceName, serviceConfig := range serviceConfigs {
		var module string
		var err error
		switch serviceProtectionLevels[serviceName] {
		case "public":
			module, err = generateServiceModule(serviceName, serviceConfig, imagesMap, "public_service.tf")
		case "private":
			module, err = generateServiceModule(serviceName, serviceConfig, imagesMap, "private_service.tf")
		}
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceName string, serviceConfig types.ServiceConfig, imagesMap map[string]string, filename string) (string, error) {
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
		"dockerImage":    imagesMap[serviceName],
		//"envVars": TODO: determine how we define env vars and then implement
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(config types.DeployConfig, imagesMap map[string]string) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range config.AppConfig.Dependencies {
		deploymentConfig := appDependency.NewAppDependency(dependency, config.AppConfig, config.AppDir, config.HomeDir).GetDeploymentConfig()
		deploymentConfig["dockerImage"] = imagesMap[dependency.Name]
		module, err := RenderTemplates(fmt.Sprintf("%s.tf", dependency.Name), deploymentConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}
