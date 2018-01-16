package terraform

import (
	"fmt"
	"path"

	"github.com/Originate/exosphere/src/types/deploy"
)

// TerraformVersion is the currently supported version of terraform
const TerraformVersion = "0.11.0"

// TerraformModulesRef is the git commit hash of the Terraform modules in Originate/exosphere we are using
const TerraformModulesRef = "89c45cc0"

// GenerateFiles generates the main terraform file given application and service configuration
func GenerateFiles(deployConfig deploy.Config) error {
	infraFileData, err := GenerateInfrastructure(deployConfig)
	if err != nil {
		return err
	}
	err = WriteToTerraformDir(infraFileData, terraformFile, deployConfig.GetInfrastructureTerraformDir())
	if err != nil {
		return err
	}
	servicesFileData, err := GenerateServices(deployConfig)
	if err != nil {
		return err
	}
	return WriteToTerraformDir(servicesFileData, terraformFile, deployConfig.GetServicesTerraformDir())
}

// GenerateCheck validates that the generated terraform file is up to date
func GenerateCheck(deployConfig deploy.Config) error {
	newInfrastructureTerraform, err := GenerateInfrastructure(deployConfig)
	if err != nil {
		return err
	}
	err = checkTerrformFile(deployConfig.GetInfrastructureTerraformDir(), newInfrastructureTerraform)
	if err != nil {
		return err
	}
	newServicesTerraform, err := GenerateServices(deployConfig)
	if err != nil {
		return err
	}
	return checkTerrformFile(deployConfig.GetServicesTerraformDir(), newServicesTerraform)
}

// checkTerrformFile checks to see if newTerraformFileContents differs from existing the terraform file
func checkTerrformFile(terraformDir, newTerraformFileContents string) error {
	currTerraformFileBytes, err := ReadTerraformFile(terraformDir)
	if err != nil {
		return err
	}
	if newTerraformFileContents != string(currTerraformFileBytes) {
		relativePath := path.Base(terraformDir)
		return fmt.Errorf("'terraform/%s/%s' is out of date. Please run 'exo generate terraform' and review the changes", relativePath, terraformFile)
	}
	return nil
}
