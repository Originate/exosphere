package template

import (
	"io/ioutil"
	"path"

	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

const applicationProjectJSONContent = `
{
  "AppName": "my-app",
  "ExocomVersion": "0.26.1",
  "AppVersion": "0.0.1",
  "AppDescription": ""
}
`

const applicationYmlContent = `name: {{AppName}}
description: {{AppDescription}}
version: {{AppVersion}}

development:
  dependencies:
    - name: exocom
      version: {{ExocomVersion}}

services:
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
	appDir := path.Join(templateDir, "template")
	if err := util.MakeDirectory(path.Join(appDir, templatesDir)); err != nil {
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
