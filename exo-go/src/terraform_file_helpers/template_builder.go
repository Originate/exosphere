package terraformFileHelpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hoisie/mustache"
)

// RenderTemplates renders templates
func RenderTemplates(templateName string, varsMap map[string]string) {
	template := getTemplate(templateName)
	data := mustache.Render(template, varsMap)
	fmt.Println(data)
	writeFile(data, "main.tf")
}

// returns a stringified template
func getTemplate(template string) string {
	data, err := Asset(fmt.Sprintf("src/terraform_file_helpers/templates/%s", template))
	if err != nil {
		log.Fatalf("Failed to read Terraform template files: %s", err)
	}
	return string(data)
}

func writeFile(data string, fileName string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %s", err)
	}

	var filePerm os.FileMode = 0744 //standard Unix file permission: rwxrw-rw-
	err = os.MkdirAll(filepath.Join(cwd, "terraform"), filePerm)
	if err != nil {
		log.Fatalf("Failed to get create directory: %s", err)
	}

	err = ioutil.WriteFile(filepath.Join(cwd, "terraform", fileName), []byte(data), filePerm)
	if err != nil {
		log.Fatalf("Failed writing Terraform files: %s", err)
	}
}
