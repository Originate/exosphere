package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-create/types"
)

func CreateApplicationYAML(appConfig types.AppConfig) error {
	const template = "name: %s\ndescription: %s\nversion: %s\n\ndependencies:\n  - name: exocom\n    version: %s\n\nservices:\n  public:\n  private:\n"
	contentBytes := []byte(fmt.Sprintf(template, appConfig.AppName, appConfig.AppDescription, appConfig.AppVersion, appConfig.ExocomVersion))
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Mkdir(path.Join(cwd, appConfig.AppName), os.FileMode(0777)); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(cwd, appConfig.AppName, "application.yml"), contentBytes, 0777)
}

func NotSet(value string) bool {
	return len(value) == 0
}
