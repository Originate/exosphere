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

// TerraformImage is the name of the terraform docker image we run commands in
const TerraformImage = "hashicorp/terraform"

// TerraformModulesRef is the git commit hash of the Terraform modules in Originate/exosphere we are using
const TerraformModulesRef = "1bf0375f"

// GenerateFile generates the main terraform file given application and service configuration
func GenerateFile(deployConfig deploy.Config) error {
	fileData, err := Generate(deployConfig)
	if err != nil {
		return err
	}
	err = WriteTerraformFile(fileData, deployConfig.TerraformDir)
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
		relativeTerraformDirPath, err := filepath.Rel(deployConfig.AppContext.Location, deployConfig.TerraformDir)
		if err != nil {
			return err
		}
		return fmt.Errorf("'%s' is out of date. Please run 'exo generate'", filepath.Join(relativeTerraformDirPath, terraformFile))
	}
	return nil
}

func generateAwsModule(deployConfig deploy.Config) (string, error) {
	varsMap := map[string]string{
		"appName":     deployConfig.AppContext.Config.Name,
		"stateBucket": deployConfig.AwsConfig.TerraformStateBucket,
		"lockTable":   deployConfig.AwsConfig.TerraformLockTable,
		"region":      deployConfig.AwsConfig.Region,
		"accountID":   deployConfig.AwsConfig.AccountID,
		"url":         deployConfig.AppContext.Config.Remote.URL,
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
		"url":                 serviceConfig.Remote.URL,
		"sslCertificateArn":   deployConfig.AwsConfig.SslCertificateArn,
		"healthCheck":         serviceConfig.Remote.HealthCheck,
		"terraformCommitHash": TerraformModulesRef,
	}
	return RenderTemplates(filename, varsMap)
}

func generateDependencyModules(deployConfig deploy.Config) (string, error) {
	dependencyModules := []string{}
	for _, dependency := range deployConfig.AppContext.Config.Remote.Dependencies {
		module, err := generateDependencyModule(dependency, deployConfig)
		if err != nil {
			return "", err
		}
		dependencyModules = append(dependencyModules, module)
	}
	for _, serviceRole := range deployConfig.AppContext.Config.GetSortedServiceRoles() {
		serviceConfig := deployConfig.AppContext.ServiceContexts[serviceRole].Config
		for _, dependency := range serviceConfig.Remote.Dependencies {
			module, err := generateDependencyModule(dependency, deployConfig)
			if err != nil {
				return "", err
			}
			dependencyModules = append(dependencyModules, module)
		}
	}
	return strings.Join(dependencyModules, "\n"), nil
}

func generateDependencyModule(dependency types.RemoteDependency, deployConfig deploy.Config) (string, error) {
	deploymentConfig, err := config.NewRemoteAppDependency(dependency, deployConfig.AppContext).GetDeploymentConfig()
	if err != nil {
		return "", err
	}
	deploymentConfig["terraformCommitHash"] = TerraformModulesRef
	return RenderTemplates(fmt.Sprintf("%s.tf", getTerraformFileName(dependency)), deploymentConfig)
}

func getTerraformFileName(dependency types.RemoteDependency) string {
	dbDependency := dependency.GetDbDependency()
	if dbDependency != "" {
		return dbDependency
	}
	return dependency.Name
}
