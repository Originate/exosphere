package serviceTemplate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

const projectJsonContent = `
{
  "ServiceType": "",
  "Description": "",
  "Author": "",
  "ServiceRole": "%s",
  "ProjectLevel": [
    "private",
    "public"
  ]
}
`

const serviceYmlContent = `type: {{ServiceType}}
description: {{Description}}
author: {{Author}}

setup:
startup:
  command: ./{{ServiceRole}}
  online-text:

messages:
  receives:
  sends:
`

func CreateTemplateDir(serviceRole string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templateDir := path.Join(cwd, "tmp")
	serviceDir := path.Join(templateDir, "template")
	if err := os.MkdirAll(serviceDir, os.FileMode(0777)); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir, serviceRole); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createServiceYMLTemplate(serviceDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create service.yml for the template")
	}
	// make a new application.yml template with service name and protection level {{}}
	return templateDir, nil
}

func RemoveTemplateDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	templateDir := path.Join(cwd, "tmp")
	return os.RemoveAll(templateDir)
}

func createProjectJSON(templateDir, serviceRole string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(fmt.Sprintf(projectJsonContent, serviceRole)), 0777)
}

func createServiceYMLTemplate(serviceDir string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(serviceYmlContent), 0777)
}
