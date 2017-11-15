package terraform

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/pkg/errors"
)

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(deployConfig types.DeployConfig) error {
	fileData, err := Generate(deployConfig)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, deployConfig.TerraformDir)
	return err
}

// Generate generates the contents of the main terraform file given application and service configuration
func Generate(deployConfig types.DeployConfig) (string, error) {
	fileData := []string{}

	moduleData, err := generateAwsModule(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateServiceModules(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate service Terraform modules")
	}
	fileData = append(fileData, moduleData)

	moduleData, err = generateDependencyModules(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)

	return strings.Join(fileData, "\n"), nil
}

func generateAwsModule(deployConfig types.DeployConfig) (string, error) {
	varsMap := map[string]string{
		"appName":     deployConfig.AppContext.Config.Name,
		"stateBucket": deployConfig.AwsConfig.TerraformStateBucket,
		"lockTable":   deployConfig.AwsConfig.TerraformLockTable,
		"region":      deployConfig.AwsConfig.Region,
		"accountID":   deployConfig.AwsConfig.AccountID,
		"url":         deployConfig.AppContext.Config.Production.URL,
		"terraformCommitHash": deployConfig.TerraformModulesRef,
	}
	return RenderTemplates("aws.tf", varsMap)
}

func generateServiceModules(deployConfig types.DeployConfig) (string, error) {
	serviceModules := []string{}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		serviceConfig := deployConfig.ServiceConfigs[serviceRole]
		module, err := generateServiceModule(serviceRole, deployConfig, serviceConfig, fmt.Sprintf("%s_service.tf", serviceConfig.Type))
		if err != nil {
			return "", err
		}
		serviceModules = append(serviceModules, module)
	}
	return strings.Join(serviceModules, "\n"), nil
}

func generateServiceModule(serviceRole string, deployConfig types.DeployConfig, serviceConfig types.ServiceConfig, filename string) (string, error) {
	varsMap := map[string]string{
		"serviceRole":         serviceRole,
		"publicPort":          serviceConfig.Production.Port,
		"cpu":                 serviceConfig.Production.CPU,
		"memory":              serviceConfig.Production.Memory,
		"url":                 serviceConfig.Production.URL,
		"sslCertificateArn":   deployConfig.AwsConfig.SslCertificateArn,
		"healthCheck":         serviceConfig.Production.HealthCheck,
		"terraformCommitHash": deployConfig.TerraformModulesRef,
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(deployConfig types.DeployConfig) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range deployConfig.AppContext.Config.Production.Dependencies {
		module, err := generateDependencyModule(dependency, deployConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		serviceConfig := deployConfig.ServiceConfigs[serviceRole]
		for _, dependency := range serviceConfig.Production.Dependencies {
			module, err := generateDependencyModule(dependency, deployConfig)
			if err != nil {
				return "", err
			}
			dependencyModules = append(dependencyModules, module)
		}
	}
	return strings.Join(dependencyModules, "\n"), nil
}

func generateDependencyModule(dependency types.ProductionDependencyConfig, deployConfig types.DeployConfig) (string, error) {
	deploymentConfig, err := config.NewAppProductionDependency(dependency, deployConfig.AppContext.Config, deployConfig.AppContext.Location).GetDeploymentConfig()
	if err != nil {
		return "", err
	}
	deploymentConfig["terraformCommitHash"] = deployConfig.TerraformModulesRef
	return RenderTemplates(fmt.Sprintf("%s.tf", getTerraformFileName(dependency)), deploymentConfig)
}

func getTerraformFileName(dependency types.ProductionDependencyConfig) string {
	dbDependency := dependency.GetDbDependency()
	if dbDependency != "" {
		return dbDependency
	}
	return dependency.Name
}
