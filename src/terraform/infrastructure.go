package terraform

import (
	"strings"

	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/pkg/errors"
)

// GenerateInfrastructure generates the contents of the main terraform file given application and service configuration
func GenerateInfrastructure(deployConfig deploy.Config) (string, error) {
	fileData := []string{}
	moduleData, err := generateAwsModule(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate AWS Terraform module")
	}
	fileData = append(fileData, moduleData)
	moduleData, err = generateDependencyModules(deployConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate application dependency Terraform modules")
	}
	fileData = append(fileData, moduleData)
	return strings.Join(fileData, "\n"), nil
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

func generateDependencyModules(deployConfig deploy.Config) (string, error) {
	dependencyModules := []string{}
	dependencies := deployConfig.AppContext.GetRemoteDependencies()
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
