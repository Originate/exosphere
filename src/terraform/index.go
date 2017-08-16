package terraform

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
)

// TerraformModulesCommitHash is the git commit hash that reflects
// which of the Terraform modules in Originate/exosphere we are using
const TerraformModulesCommitHash = "8786f912"

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(deployConfig types.DeployConfig, imagesMap map[string]string) error {
	fileData, err := Generate(deployConfig, imagesMap)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, deployConfig.TerraformDir)
	return err
}

// Generate generates the contents of the main terraform file given application and service configuration
func Generate(deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	fileData := []string{}

	moduleData, err := generateAwsModule(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	serviceProtectionLevels := deployConfig.AppConfig.GetServiceProtectionLevels()
	moduleData, err = generateServiceModules(deployConfig.ServiceConfigs, serviceProtectionLevels, imagesMap)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(deployConfig, imagesMap)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	return strings.Join(fileData, "\n"), nil
}

func generateAwsModule(deployConfig types.DeployConfig) (string, error) {
	varsMap := map[string]string{
		"appName":     deployConfig.AppConfig.Name,
		"stateBucket": deployConfig.AwsConfig.TerraformStateBucket,
		"lockTable":   deployConfig.AwsConfig.TerraformLockTable,
		"region":      deployConfig.AwsConfig.Region,
		"url":         deployConfig.AppConfig.Production["url"],
		"terraformCommitHash": TerraformModulesCommitHash,
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
		"serviceRole":         serviceName,
		"startupCommand":      string(command),
		"publicPort":          serviceConfig.Production["public-port"],
		"cpu":                 serviceConfig.Production["cpu"],
		"memory":              serviceConfig.Production["memory"],
		"url":                 serviceConfig.Production["url"],
		"healthCheck":         serviceConfig.Production["health-check"],
		"dockerImage":         imagesMap[serviceName],
		"terraformCommitHash": TerraformModulesCommitHash,
		//"envVars": TODO: determine how we define env vars and then implement
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range deployConfig.AppConfig.Dependencies {
		deploymentConfig, err := config.NewAppDependency(dependency, deployConfig.AppConfig, deployConfig.AppDir, deployConfig.HomeDir).GetDeploymentConfig()
		if err != nil {
			return "", err
		}
		deploymentConfig["dockerImage"] = imagesMap[dependency.Name]
		deploymentConfig["terraformCommitHash"] = TerraformModulesCommitHash
		module, err := RenderTemplates(fmt.Sprintf("%s.tf", dependency.Name), deploymentConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}
