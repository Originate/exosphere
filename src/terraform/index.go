package terraform

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
)

// terraformModulesCommitHash is the git commit hash that reflects
// which of the Terraform modules in Originate/exosphere we are using
const terraformModulesCommitHash = "fa599050"

// dependenciesMap is a map of dependency name to the Terraform template
// name it uses. It's for managing multiple dependencies with different
// names but use the same Terraform template
const dependenciesMap = map[string]string{
	"postgres": "rds",
}

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
	moduleData, err = generateServiceModules(deployConfig, serviceProtectionLevels, imagesMap)
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
		"awsProfile":  deployConfig.AwsConfig.Profile,
		"stateBucket": deployConfig.AwsConfig.TerraformStateBucket,
		"lockTable":   deployConfig.AwsConfig.TerraformLockTable,
		"region":      deployConfig.AwsConfig.Region,
		"accountID":   deployConfig.AwsConfig.AccountID,
		"url":         deployConfig.AppConfig.Production.URL,
		"terraformCommitHash": terraformModulesCommitHash,
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(deployConfig types.DeployConfig, serviceProtectionLevels, imagesMap map[string]string) (string, error) {
	serviceModules := []string{}
	for serviceName, serviceConfig := range deployConfig.ServiceConfigs {
		module, err := generateServiceModule(serviceName, deployConfig.AwsConfig, serviceConfig, imagesMap, fmt.Sprintf("%s_service.tf", serviceProtectionLevels[serviceName]))
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceName string, awsConfig types.AwsConfig, serviceConfig types.ServiceConfig, imagesMap map[string]string, filename string) (string, error) {
	varsMap := map[string]string{
		"serviceRole":         serviceName,
		"publicPort":          serviceConfig.Production.PublicPort,
		"cpu":                 serviceConfig.Production.CPU,
		"memory":              serviceConfig.Production.Memory,
		"url":                 serviceConfig.Production.URL,
		"sslCertificateArn":   awsConfig.SslCertificateArn,
		"healthCheck":         serviceConfig.Production.HealthCheck,
		"dockerImage":         imagesMap[serviceName],
		"terraformCommitHash": terraformModulesCommitHash,
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range deployConfig.AppConfig.Production.Dependencies {
		module, err := generateDependencyModule(dependency, deployConfig, imagesMap)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	for _, serviceConfigs := range deployConfig.ServiceConfigs {
		for _, dependency := range serviceConfigs.Production.Dependencies {
			module, err := generateDependencyModule(dependency, deployConfig, imagesMap)
			if err != nil {
				return "", err
			}
			dependencyModules = append(dependencyModules, module)
		}
	}
	return strings.Join(dependencyModules, "\n"), nil
}

func generateDependencyModule(dependency types.DependencyConfig, deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	deploymentConfig, err := config.NewAppDependency(dependency, deployConfig.AppConfig, deployConfig.AppDir, deployConfig.HomeDir).GetDeploymentConfig()
	if err != nil {
		return "", err
	}
	deploymentConfig["dockerImage"] = imagesMap[dependency.Name]
	deploymentConfig["terraformCommitHash"] = terraformModulesCommitHash
	dependencyName := dependency.Name
	if dependenciesMap[dependency.Name] != "" {
		dependencyName = dependenciesMap[dependency.Name]
	}
	return RenderTemplates(fmt.Sprintf("%s.tf", dependencyName), deploymentConfig)
}
