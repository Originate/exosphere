package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/terraform/remotedependencies"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/util"
	"github.com/hoisie/mustache"
	"github.com/pkg/errors"
)

const terraformFile = "main.tf"

// RenderTemplates renders a Terraform template
func RenderTemplates(templateName string, varsMap map[string]string) (string, error) {
	template, err := getTemplate(templateName)
	if err != nil {
		return "", err
	}
	return mustache.Render(template, varsMap), nil
}

// RenderRemoteTemplates renders a Terraform template
func RenderRemoteTemplates(dependencyType string, templateConfig map[string]string) (string, error) {
	template, err := getRemoteTemplate(dependencyType)
	if err != nil {
		return "", err
	}
	return mustache.Render(template, templateConfig), nil
}

// WriteTerraformFile writes the main Terraform file to the given path
func WriteTerraformFile(data string, terraformDir string) error {
	err := util.MakeDirectory(terraformDir)
	if err != nil {
		return errors.Wrap(err, "Failed to get create 'terraform' directory")
	}

	err = writeGitIgnore(terraformDir)
	if err != nil {
		return errors.Wrap(err, "Failed to write 'terraform/.gitignore' file")
	}

	var filePerm os.FileMode = 0744 //standard Unix file permission: rwxrw-rw-
	err = ioutil.WriteFile(filepath.Join(terraformDir, terraformFile), []byte(data), filePerm)
	if err != nil {
		return errors.Wrap(err, "Failed to write 'terraform/main.tf' file")
	}
	return nil
}

func writeGitIgnore(terraformDir string) error {
	gitIgnore := `.terraform/
terraform.tfstate
terraform.tfstate.backup`
	gitIgnorePath := filepath.Join(terraformDir, ".gitignore")
	fileExists, err := util.DoesFileExist(gitIgnorePath)
	if !fileExists {
		return ioutil.WriteFile(gitIgnorePath, []byte(gitIgnore), 0744)
	}
	return err
}

func getTemplate(template string) (string, error) {
	data, err := remotedependencies.Asset(fmt.Sprintf("src/terraform/templates/%s", template))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}

func getRemoteTemplate(dependencyType string) (string, error) {
	data, err := remotedependencies.Asset(fmt.Sprintf("remote-dependency-templates/%s/dependency.tf", dependencyType))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}

// ReadTerraformFile reads the contents of the main terraform file
func ReadTerraformFile(deployConfig deploy.Config) ([]byte, error) {
	terraformFilePath := filepath.Join(deployConfig.TerraformDir, terraformFile)
	fileExists, err := util.DoesFileExist(terraformFilePath)
	if fileExists {
		return ioutil.ReadFile(terraformFilePath)
	}
	return []byte{}, err
}
