package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/assets"
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

// WriteToNestedTerraformDir writes data to the terraform dir
func WriteToNestedTerraformDir(data, fileName, nestedTerraformDir string) error {
	err := util.MakeDirectory(nestedTerraformDir)
	if err != nil {
		return err
	}
	err = writeGitIgnore(filepath.Dir(nestedTerraformDir))
	if err != nil {
		return err
	}
	var filePerm os.FileMode = 0744 //standard Unix file permission: rwxrw-rw-
	return ioutil.WriteFile(filepath.Join(nestedTerraformDir, fileName), []byte(data), filePerm)
}

func writeGitIgnore(terraformDir string) error {
	gitIgnore := `.terraform/
*.tfvars
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
	data, err := assets.Asset(fmt.Sprintf("src/terraform/templates/%s", template))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}

func getRemoteTemplate(dependencyType string) (string, error) {
	data, err := assets.Asset(fmt.Sprintf("remote-dependency-templates/%s/dependency.tf", dependencyType))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}

// ReadTerraformFile reads the contents of the main terraform file
func ReadTerraformFile(terraformDir string) ([]byte, error) {
	terraformFilePath := filepath.Join(terraformDir, terraformFile)
	fileExists, err := util.DoesFileExist(terraformFilePath)
	if fileExists {
		return ioutil.ReadFile(terraformFilePath)
	}
	return []byte{}, err
}
