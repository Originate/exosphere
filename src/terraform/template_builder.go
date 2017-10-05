package terraform

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Originate/exosphere/src/types"
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

// WriteTerraformFile writes the main Terraform file to the given path
func WriteTerraformFile(data string, terraformDir string) error {
	var filePerm os.FileMode = 0744 //standard Unix file permission: rwxrw-rw-
	err := os.MkdirAll(terraformDir, filePerm)
	if err != nil {
		return errors.Wrap(err, "Failed to get create directory")
	}

	err = ioutil.WriteFile(filepath.Join(terraformDir, terraformFile), []byte(data), filePerm)
	if err != nil {
		return errors.Wrap(err, "Failed writing Terraform files")
	}
	return nil
}

func getTemplate(template string) (string, error) {
	data, err := Asset(fmt.Sprintf("src/terraform/templates/%s", template))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}

// ReadTerraformFile reads the contents of the main terraform file
func ReadTerraformFile(deployConfig types.DeployConfig) ([]byte, error) {
	terraformFilePath := filepath.Join(deployConfig.TerraformDir, terraformFile)
	fileExists, err := os.Stat(terraformFilePath)
	if fileExists {
		return ioutil.ReadFile(terraformFilePath)
	}
	return []byte{}, err
}

// CheckTerraformFile makes sure that the generated terraform file hasn't changed from the previous one
// It returns an error if they differ. Used for deploying from CI servers
func CheckTerraformFile(deployConfig types.DeployConfig, prevTerraformFileContents []byte) error {
	terraformFilePath := filepath.Join(deployConfig.TerraformDir, terraformFile)
	generatedTerraformFileContents, err := ioutil.ReadFile(terraformFilePath)
	if err != nil {
		return err
	}
	if !bytes.Equal(generatedTerraformFileContents, prevTerraformFileContents) {
		return errors.New("'terraform/main.tf' file has changed. Please deploy manually to review these changes")
	}
	return nil
}
