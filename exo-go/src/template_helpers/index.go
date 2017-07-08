package templateHelpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/process_helpers"
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

// isValidTemplateDir returns true if the directory templateDir is a valid
// boilr template directory
func isValidTemplateDir(templateDir string) bool {
	return osHelpers.FileExists(path.Join(templateDir, "project.json")) && osHelpers.DirectoryExists(path.Join(templateDir, "template"))
}

// CreateServiceYML creates service.yml for the service serviceRole by creating
// a boilr template for service.yml, making boilr do the scaffolding and finally
// removing the template
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
		log.Fatalf("Failed to remove service.yml template: %s", err)
	}
}

// CreateApplicationTemplateDir creates a temporary boilr template directory
// for the application
func CreateApplicationTemplateDir() (string, error) {
	templateDir, err := ioutil.TempDir("", "application")
	if err != nil {
		return templateDir, errors.Wrap(err, "Failed to create temp dir for application template")
	}
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

// AddTemplate fetches a remote template from GitHub and stores it
// under templateDir, returns an error if any
func AddTemplate(gitURL, templateName, templateDir string) error {
	if output, err := processHelpers.Run("", "git", "submodule", "add", "--name", templateName, gitURL, templateDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to fetch template from %s: %s\n", gitURL, output))
	}
	return nil
}

// FetchTemplates fetches updates for all existing remote templates
func FetchTemplates() error {
	if output, err := processHelpers.Run("", "git", "submodule", "foreach", "git", "pull", "origin", "master"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to fetch updates for existing templates: %s\n", output))
	}
	return nil
}

// RemoveTemplate removes the given template from the application
func RemoveTemplate(templateName, templateDir string) error {
	if output, err := processHelpers.Run("", "git", "submodule", "deinit", "-f", templateDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to deinit the template submodule: %s\n", output))
	}
	if output, err := processHelpers.Run("", "rm", "-rf", fmt.Sprintf(".git/modules/%s", templateName)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to force remove the template submodule: %s\n", output))
	}
	if output, err := processHelpers.Run("", "git", "rm", "-f", templateDir); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to git remove the template submodule: %s\n", output))
	}
	return nil
}

// UpdateTemplate updates remote template according to the given GitHub URL,
// return an error if any
func UpdateTemplate(gitURL, templateName string) error {
	if output, err := processHelpers.Run("", "git", "config", "--file=.gitmodules", fmt.Sprintf("submodule.%s.url", templateName), gitURL); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to update git URL for template %s: %s\n", templateName, output))
	}
	if output, err := processHelpers.Run("", "git", "submodule", "sync"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to sync the submodule %s: %s\n", templateName, output))
	}
	if output, err := processHelpers.Run("", "git", "submodule", "update", "--init", "--recursive", "--remote"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to update the submodule %s: %s\n", templateName, output))
	}
	return nil
}
