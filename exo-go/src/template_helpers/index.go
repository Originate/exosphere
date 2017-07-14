package templateHelpers

import (
	"fmt"
	"io/ioutil"
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

const templatesDir = ".exosphere"

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
func CreateServiceYML(serviceRole string) error {
	templateDir, err := createServiceTemplateDir(serviceRole)
	if err != nil {
		return err
	}
	serviceYmlTemplate, err := template.Get(templateDir)
	if err != nil {
		return err
	}
	if err = serviceYmlTemplate.Execute(serviceRole); err != nil {
		return err
	}
	if err = os.RemoveAll(templateDir); err != nil {
		return err
	}
	return nil
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

// GetTemplates returns a slice of all template names found in the ".exosphere"
// folder of the application
func GetTemplates() (result []string, err error) {
	subdirectories, err := osHelpers.GetSubdirectories(templatesDir)
	if err != nil {
		return result, err
	}
	for _, directory := range subdirectories {
		if isValidTemplateDir(path.Join(templatesDir, directory)) {
			result = append(result, directory)
		}
	}
	return result, nil
}

// HasTemplateDirectory returns whether or not there is an ".exosphere" folder
func HasTemplateDirectory() bool {
	return osHelpers.DirectoryExists(templatesDir) && !osHelpers.IsEmpty(templatesDir)
}

// CreateTmpServiceDir makes bolir scaffold the template chosenTemplate
// and store the scaffoled service folder in a tmp folder, and finally
// returns the path to the tmp folder
func CreateTmpServiceDir(chosenTemplate string) (string, error) {
	templateDir := path.Join(templatesDir, chosenTemplate)
	template, err := template.Get(templateDir)
	if err != nil {
		return "", err
	}
	serviceTmpDir, err := ioutil.TempDir("", "service-tmp")
	if err != nil {
		return "", err
	}
	if err = template.Execute(serviceTmpDir); err != nil {
		return "", err
	}
	return serviceTmpDir, nil
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
	denitCommand := []string{"git", "submodule", "deinit", "-f", templateDir}
	removeModulesCommand := []string{"rm", "-rf", fmt.Sprintf(".git/modules/%s", templateName)}
	gitRemoveCommand := []string{"git", "rm", "-f", templateDir}
	return processHelpers.RunSeries("", [][]string{denitCommand, removeModulesCommand, gitRemoveCommand})
}
