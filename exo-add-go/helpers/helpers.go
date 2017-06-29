package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-add-go/types"
	"github.com/pkg/errors"
	"github.com/tmrts/boilr/pkg/template"
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

func GetTemplateDirs() []string {
	templateDir := ".exosphere"
	if !directoryExists(templateDir) || isEmpty(templateDir) {
		log.Fatal("no templates found")
	}
	templates := []string{}
	for _, directory := range getSubdirectories(templateDir) {
		fmt.Println("directory", directory)
		if isValidTemplateDir(path.Join(templateDir, directory)) {
			templates = append(templates, directory)
		}
	}
	return templates
}

func isValidTemplateDir(templateDir string) bool {
	fmt.Println("templateDir", templateDir)
	return fileExists(path.Join(templateDir, "project.json")) && directoryExists(path.Join(templateDir, "template"))
}

func createProjectJSON(templateDir, serviceDirectory string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(fmt.Sprintf(projectJsonContent, serviceDirectory)), 0777)
}

func createServiceYMLTemplate(serviceDir string) error {
	return ioutil.WriteFile(path.Join(serviceDir, "service.yml"), []byte(serviceYmlContent), 0777)
}

func createTemplateDir(serviceDirectory string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templateDir := path.Join(cwd, "tmp")
	serviceDir := path.Join(templateDir, "template")
	if err := os.MkdirAll(serviceDir, os.FileMode(0777)); err != nil {
		return templateDir, errors.Wrap(err, "Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir, serviceDirectory); err != nil {
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
	fmt.Println("making service.yml")
	templatePath, err := createTemplateDir(serviceDirectory)
	fmt.Println("created the template service.yml")
	if err != nil {
		log.Fatalf("Failed to create the template: %s", err)
	}
	template, err := template.Get(templatePath)
	fmt.Println("got the template service.yml")
	if err != nil {
		log.Fatalf("Failed to fetch service.yml template: %s", err)
	}
	fmt.Println("executing service.yml")
	if err = template.Execute(serviceDirectory); err != nil {
		log.Fatalf("Failed to create service.yml: %s", err)
	}
	fmt.Println("done executing service.yml")
	if err = removeTemplateDir(); err != nil {
		log.Fatalf("Failed to remove the template: %s", err)
	}
	fmt.Println("removed")
}

func fileExists(filePath string) bool {
	// _, err := os.Stat(filePath)
	return true
}

func directoryExists(dirPath string) bool {
	f, _ := os.Stat(dirPath)
	return f.IsDir()
}

func isDirectory(dirPath string) bool {
	f, _ := os.Stat(dirPath)
	return f.IsDir()
}

func isEmpty(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	_, err = f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
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
