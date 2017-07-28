package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

const serviceProjectJSONContent = `
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

func createServiceYMLTemplate(serviceDir, serviceRole string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(fmt.Sprintf(serviceYmlContent, serviceRole)), 0777)
}

// CreateServiceTemplateDir creates a temporary boilr template directory
// for the service
func CreateServiceTemplateDir(serviceRole string) (string, error) {
	templateDir, err := ioutil.TempDir("", "service-yml")
	if err != nil {
		return templateDir, errors.Wrap(err, "Failed to create temp dir for service.yml template")
	}
	serviceYMLDir := path.Join(templateDir, "template")
	if err := os.Mkdir(serviceYMLDir, 0700); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir, serviceProjectJSONContent); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createServiceYMLTemplate(serviceYMLDir, serviceRole); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create service.yml for the template")
	}
	return templateDir, nil
}
