package templateHelpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/template"
)

const applicationProjectJSONContent = `
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

func createProjectJSON(templateDir string, content string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(content), 0777)
}

func createApplicationYML(appDir string) error {
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(applicationYmlContent), 0777)
}

func createServiceYMLTemplate(serviceDir, serviceRole string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(fmt.Sprintf(serviceYmlContent, serviceRole)), 0777)
}

func createServiceTemplateDir(serviceRole string) (string, error) {
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

// CreateServiceYML creates service.yml for the service serviceRole by creating
// a boilr template for service.yml, making boilr do the scaffolding and finally
// remove the template
func CreateServiceYML(serviceRole string) {
	templateDir, err := createServiceTemplateDir(serviceRole)
	if err != nil {
		log.Fatalf("Failed to create the service.yml template: %s", err)
	}
	serviceYmlTemplate, err := template.Get(templateDir)
	if err != nil {
		log.Fatalf("Failed to fetch service.yml template: %s", err)
	}
	if err = serviceYmlTemplate.Execute(serviceRole); err != nil {
		log.Fatalf("Failed to create service.yml: %s", err)
	}
	if err = os.RemoveAll(templateDir); err != nil {
		log.Fatalf("Failed to remove the template: %s", err)
	}
}

// CreateTemplateDir creates and populates a temporary template directory
// to be used with boilr
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
	if err := createProjectJSON(templateDir, applicationProjectJSONContent); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createApplicationYML(appDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create application.yml for the template")
	}
	return templateDir, nil
}

// GetTemplates returns a slice of all template names found in the ".exosphere"
// folder of the application
func GetTemplates() []string {
	templatesDir := ".exosphere"
	if !osHelpers.DirectoryExists(templatesDir) || osHelpers.IsEmpty(templatesDir) {
		fmt.Println("no templates found\n\nPlease add templates to the \".exosphere\" folder of your code base.")
		os.Exit(1)
	}
	templates := []string{}
	for _, directory := range osHelpers.GetSubdirectories(templatesDir) {
		if isValidTemplateDir(path.Join(templatesDir, directory)) {
			templates = append(templates, directory)
		}
	}
	return templates
}

// IsValidTemplateDir returns true if the directory templateDir is a valid
// boilr template directory
func isValidTemplateDir(templateDir string) bool {
	return osHelpers.FileExists(path.Join(templateDir, "project.json")) && osHelpers.DirectoryExists(path.Join(templateDir, "template"))
}

// CreateTmpServiceDir makes bolir scaffold the template chosenTemplate
// and store the scaffoled service folder in a tmp folder, and finally
// returns the path to the tmp folder
func CreateTmpServiceDir(chosenTemplate string) string {
	templateDir := path.Join(".exosphere", chosenTemplate)
	serviceTemplate, err := template.Get(templateDir)
	if err != nil {
		log.Fatalf("Failed to fetch %s template: %s", chosenTemplate, err)
	}
	serviceTmpDir, err := ioutil.TempDir("", "service-tmp")
	if err != nil {
		log.Fatalf(`Failed to create a tmp folder for the service "%s": %s`, chosenTemplate, err)
	}
	if err = serviceTemplate.Execute(serviceTmpDir); err != nil {
		log.Fatalf(`Failed to create the service "%s": %s`, chosenTemplate, err)
	}
	return serviceTmpDir
}

// RemoveTemplateDir removes the temporary template directory
func RemoveTemplateDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	templateDir := path.Join(cwd, "tmp")
	return os.RemoveAll(templateDir)
}
