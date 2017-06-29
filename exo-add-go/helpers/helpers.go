package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add-go/types"
)

const projectJsonContent = `
{
  "ServiceType": "",
  "Description": "",
  "Author": "",
  "ServiceRole": %s
  "ProjectLevel": [
  	"private",
  	"public"
  ]
}
`

const serviceYmlContent = `type: {{ServiceType}}
description: {{Description}}
author: {{Author}}

setup: echo "installing dependencies for 'crasher' service"
startup:
  command: ./{{ServiceRole}}
  online-text:

messages:
  receives:
  sends:
`

func getTemplateDirs() []string {
	if !directoryExists(".exosphere") {
		log.Fatal("no templates found")
	}
	templates := []string{}
	for _, directory := range getSubdirectories(".exosphere") {
		if isValidTemplateDir() {
			templates = append(templates, directory)
		}
	}
	return templates
}

func isValidTemplateDir(templateDir string) {
	return fileExists(path.Join(templateDir, "project.json")) && directoryExists(path.Join(templateDir, "template"))
}

func createProjectJSON(templateDir string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(projectJsonContent), 0777)
}

func createServiceYMLTemplate(serviceDir string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(applicationYmlContent), 0777)
}

func createTemplateDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templateDir := path.Join(cwd, "tmp")
	serviceDir := path.Join(templateDir, "template")
	if err := os.MkdirAll(serviceDir), os.FileMode(0777)); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create project.json for the template")
	}
	if err := createServiceYMLTemplate(serviceDir); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create service.yml for the template")
	}
	// make a new application.yml template with service name and protection level {{}}
	return templateDir, nil
}

func removeTemplateDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	templateDir := path.Join(cwd, "tmp")
	return os.RemoveAll(templateDir)
}

func CreateServiceYML(serviceDirectory string) {
	templatePath, err := helpers.createTemplateDir()
	if err != nil {
		log.Fatalf("Failed to create the template: %s", err)
	}
	template, err := template.Get(templatePath)
	if err != nil {
		log.Fatalf("Failed to fetch service.yml template: %s", err)
	}
	if err = template.Execute(serviceDirectory); err != nil {
		log.Fatalf("Failed to create service.yml: %s", err)
	}
	if err = helpers.removeTemplateDir(); err != nil {
		log.Fatalf("Failed to remove the template: %s", err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsExist(err)
}

func directoryExists(dirPath string) {
	f, err := os.Stat(filePath)
	return os.IsExist(err) && f.IsDir()
}


func contains(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

func CheckForService(serviceRole string, existingServices []string) {
	if contains(existingServices, serviceRole) {
		fmt.Printf("Service %v already exists in this application\n", serviceRole)
		os.Exit(1)
	}
}

func GetExistingServices(services map[string]map[string]types.Service) []string {
	existingServices := []string{}
	for _, serviceConfigs := range services {
		for service := range serviceConfigs {
			existingServices = append(existingServices, service)
		}
	}
	return existingServices
}

func getSubdirectories(directory string) []string {
	subDirectories := []string{}
	entries, _ := ioutil.ReadDir(directory)
  for _, entry := range entries {
  	if isDirectory(path.Join(directory, entry.Name())) {
  		subDirectories = append(subDirectories, entry.Name())
  	}
  }
  return subDirectories
}
