package terraformFileHelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hoisie/mustache"
	"github.com/pkg/errors"
)

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

	err = ioutil.WriteFile(filepath.Join(terraformDir, "main.tf"), []byte(data), filePerm)
	if err != nil {
		return errors.Wrap(err, "Failed writing Terraform files")
	}
	return nil
}

func getTemplate(template string) (string, error) {
	data, err := Asset(fmt.Sprintf("src/terraform_file_helpers/templates/%s", template))
	if err != nil {
		return "", errors.Wrap(err, "Failed to read Terraform template files")
	}
	return string(data), nil
}
