package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func createProjectJSON(templateDir string) error {
	content := "{\n\"AppName\": \"my-app\",\n\"ExocomVersion\": \"0.22.1\",\n\"AppVersion\": \"0.0.1\",\n\"AppDescription\": \"\"\n}"
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(content), 0777)
}

func createApplicationYAML(appDir string) error {
	content := "name: {{AppName}}\ndescription: {{AppDescription}}\nversion: {{AppVersion}}\n\ndependencies:\n  - name: exocom\n    version: {{ExocomVersion}}\n\nservices:\n  public:\n  private:\n"
	return ioutil.WriteFile(path.Join(appDir, "application.yml"), []byte(content), 0777)
}

func CreateTemplate() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	templateDir := path.Join(cwd, "tmp")
	appDir := path.Join(templateDir, "template/{{AppName}}")
	if err := os.MkdirAll(path.Join(appDir, ".exosphere"), os.FileMode(0777)); err != nil {
		return templateDir, fmt.Errorf("Failed to create the neccessary directories for the template")
	}
	if err := createProjectJSON(templateDir); err != nil {
		return templateDir, fmt.Errorf("Failed to create project.json for the template")
	}
	if err := createApplicationYAML(appDir); err != nil {
		return templateDir, fmt.Errorf("Failed to create application.yml for the template")
	}
	return templateDir, nil
}

func RemoveTemplate() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	templateDir := path.Join(cwd, "tmp")
	return os.RemoveAll(templateDir)
}
