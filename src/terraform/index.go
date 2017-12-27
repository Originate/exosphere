package terraform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/pkg/errors"
)

// TerraformVersion is the currently supported version of terraform
const TerraformVersion = "0.11.0"

// TerraformModulesRef is the git commit hash of the Terraform modules in Originate/exosphere we are using
const TerraformModulesRef = "30894145"

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(deployConfig deploy.Config) error {
	fileData, err := Generate(deployConfig)
	if err != nil {
		return err
	}
	err = WriteToTerraformDir(fileData, terraformFile, deployConfig.GetTerraformDir())
	return err
}

// Generate generates the contents of the main terraform file given application and service configuration
func Generate(deployConfig deploy.Config) (string, error) {
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

// GenerateCheck validates that the generated terraform file is up to date
func GenerateCheck(deployConfig deploy.Config) error {
	currTerraformFileBytes, err := ReadTerraformFile(deployConfig)
	if err != nil {
		return err
	}
	newTerraformFileContents, err := Generate(deployConfig)
	if err != nil {
		return err
	}
	if newTerraformFileContents != string(currTerraformFileBytes) {
		return fmt.Errorf("'%s' is out of date. Please run 'exo generate terraform' and review the changes", filepath.Join(deployConfig.GetRelativeTerraformDir(), terraformFile))
	}
	return nil
}

func generateAwsModule(deployConfig deploy.Config) (string, error) {
	varsMap := map[string]string{
		"appName":             deployConfig.AppContext.Config.Name,
		"lockTable":           deployConfig.AwsConfig.TerraformLockTable,
		"terraformCommitHash": TerraformModulesRef,
		"terraformVersion":    TerraformVersion,
	}
	return RenderTemplates("aws.tf", varsMap)
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

func generateDependencyModules(deployConfig deploy.Config) (string, error) {
	dependencyModules := []string{}
	dependencies := config.GetAllRemoteDependencies(deployConfig.AppContext)
	for _, dependencyName := range deployConfig.AppContext.GetSortedRemoteDependencyNames() {
		module, err := generateDependencyModule(dependencyName, dependencies[dependencyName], deployConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	return strings.Join(dependencyModules, "\n"), nil
}

func generateDependencyModule(dependencyName string, dependency types.RemoteDependency, deployConfig deploy.Config) (string, error) {
	templateConfig := dependency.TemplateConfig
	templateConfig["terraformCommitHash"] = TerraformModulesRef
	return RenderRemoteTemplates(dependency.Type, templateConfig)
}
