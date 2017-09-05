package template

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

const applicationProjectJSONContent = `
{
  "AppName": "my-app",
  "ExocomVersion": "0.24.0",
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
  worker:
`

func createApplicationYML(appDir string) error {
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(applicationYmlContent), 0777)
}

// CreateApplicationTemplateDir creates a temporary boilr template directory
// for the application
func CreateApplicationTemplateDir() (string, error) {
	templateDir, err := ioutil.TempDir("", "application")
	if err != nil {
		return templateDir, errors.Wrap(err, "Failed to create temp dir for application template")
	}
	appDir := path.Join(templateDir, "template/{{AppName}}")
	if err := os.MkdirAll(path.Join(appDir, templatesDir), os.FileMode(0777)); err != nil {
		return templateDir, err
	}
	if err := createProjectJSON(templateDir, applicationProjectJSONContent); err != nil {
		return templateDir, err
	}
	if err := createApplicationYML(appDir); err != nil {
		return templateDir, err
	}
	return templateDir, nil
}
