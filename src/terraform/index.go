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
const terraformModulesCommitHash = "1bb2c93b"

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
	moduleData, err = generateServiceModules(deployConfig, serviceProtectionLevels)
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
		"accountID":   deployConfig.AwsConfig.AccountID,
		"url":         deployConfig.AppConfig.Production.URL,
		"terraformCommitHash": terraformModulesCommitHash,
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(deployConfig types.DeployConfig, serviceProtectionLevels map[string]string) (string, error) {
	serviceModules := []string{}
	for _, serviceName := range deployConfig.AppConfig.GetServiceNames() {
		serviceConfig := deployConfig.ServiceConfigs[serviceName]
		module, err := generateServiceModule(serviceName, deployConfig.AwsConfig, serviceConfig, fmt.Sprintf("%s_service.tf", serviceProtectionLevels[serviceName]))
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceName string, awsConfig types.AwsConfig, serviceConfig types.ServiceConfig, filename string) (string, error) {
	varsMap := map[string]string{
		"serviceRole":         serviceName,
		"publicPort":          serviceConfig.Production.PublicPort,
		"cpu":                 serviceConfig.Production.CPU,
		"memory":              serviceConfig.Production.Memory,
		"url":                 serviceConfig.Production.URL,
		"sslCertificateArn":   awsConfig.SslCertificateArn,
		"healthCheck":         serviceConfig.Production.HealthCheck,
		"terraformCommitHash": terraformModulesCommitHash,
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	dependencyModules := []string{}
	generatedDependencies := map[string]string{}
	for _, dependency := range deployConfig.AppConfig.Production.Dependencies {
		module, err := generateDependencyModule(dependency, deployConfig, imagesMap)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
		generatedDependencies[dependency.Name] = dependency.Version
	}
	for _, serviceName := range deployConfig.AppConfig.GetServiceNames() {
		serviceConfig := deployConfig.ServiceConfigs[serviceName]
		for _, dependency := range serviceConfig.Production.Dependencies {
			if generatedDependencies[dependency.Name] == "" {
				module, err := generateDependencyModule(dependency, deployConfig, imagesMap)
				if err != nil {
					return "", err
				}
				dependencyModules = append(dependencyModules, module)
			}
		}
	}
	return strings.Join(dependencyModules, "\n"), nil
}

func generateDependencyModule(dependency types.ProductionDependencyConfig, deployConfig types.DeployConfig, imagesMap map[string]string) (string, error) {
	deploymentConfig, err := config.NewAppProductionDependency(dependency, deployConfig.AppConfig, deployConfig.AppDir).GetDeploymentConfig()
	if err != nil {
		return "", err
	}
	deploymentConfig["dockerImage"] = imagesMap[dependency.Name]
	deploymentConfig["terraformCommitHash"] = terraformModulesCommitHash
	return RenderTemplates(fmt.Sprintf("%s.tf", getTerraformFileName(dependency)), deploymentConfig)
}

func getTerraformFileName(dependency types.ProductionDependencyConfig) string {
	dbDependency := dependency.GetDbDependency()
	if dbDependency != "" {
		return dbDependency
	}
	return dependency.Name
}
