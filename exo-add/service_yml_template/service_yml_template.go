package serviceYmlTemplate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/template"
)

const projectJsonContent = `
{
  "ServiceType": "",
  "Description": "",
  "Author": ""
}
`

const serviceYmlContent = `type: {{ServiceType}}
description: {{Description}}
author: {{Author}}

setup:
startup:
  command: ./%s
  online-text:

messages:
  receives:
  sends:
`

// CreateServiceYML creates service.yml for the service serviceRole by creating
// a boilr template for service.yml, making boilr do the scaffolding and finally
// remove the template
func CreateServiceYML(serviceRole string) {
	templateDir, err := createTemplateDir(serviceRole)
	if err != nil {
		log.Fatalf("Failed to create the service.yml template: %s", err)
	}
	template, err := template.Get(templateDir)
	if err != nil {
		log.Fatalf("Failed to fetch service.yml template: %s", err)
	}
	if err = template.Execute(serviceRole); err != nil {
		log.Fatalf("Failed to create service.yml: %s", err)
	}
	if err = os.RemoveAll(templateDir); err != nil {
		log.Fatalf("Failed to remove the template: %s", err)
	}
}

func createTemplateDir(serviceRole string) (string, error) {
	templateDir, err := ioutil.TempDir("", "service-yml")
	if err != nil {
		return templateDir, errors.Wrap(err, "Failed to create temp dir for service.yml template")
	}
	serviceYMLDir := path.Join(templateDir, "template")
	if err := os.Mkdir(serviceYMLDir, 0777); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createServiceYMLTemplate(serviceYMLDir, serviceRole); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create service.yml for the template")
	}
	return templateDir, nil
}

func createProjectJSON(templateDir string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(projectJsonContent), 0777)
}

func createServiceYMLTemplate(serviceDir, serviceRole string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(fmt.Sprintf(serviceYmlContent, serviceRole)), 0777)
}
