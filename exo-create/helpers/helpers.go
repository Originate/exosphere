package helpers

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

const projectJsonContent = `
{
  "AppName": "my-app",
  "ExocomVersion": "0.22.1",
  "AppVersion": "0.0.1",
  "AppDescription": ""
}
`

const applicationYmlContent = `name: {{AppName}}
description: {{AppDescription}}
version: {{AppVersion}}

dependencies:
  - name: exocom
    version: {{ExocomVersion}}

services:
  public:
  private:
`

func createProjectJSON(templateDir string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(projectJsonContent), 0777)
}

func createApplicationYML(appDir string) error {
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(applicationYmlContent), 0777)
}

func CreateTemplateDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templateDir := path.Join(cwd, "tmp")
	appDir := path.Join(templateDir, "template/{{AppName}}")
	if err := os.MkdirAll(path.Join(appDir, ".exosphere"), os.FileMode(0777)); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createApplicationYML(appDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create application.yml for the template")
	}
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
