package helpers

import (
	"html/template"
	"os"
	"path"

	"github.com/Originate/exosphere/exo-create/types"
)

func CreateApplicationYAML(appConfig types.AppConfig) error {
	dir, err := os.Executable()
	if err != nil {
		return err
	}
	templatePath := path.Join(dir, "..", "..", "..", "exosphere-shared", "templates", "create-app", "application-go.yml")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Mkdir(path.Join(cwd, appConfig.AppName), os.FileMode(0777)); err != nil {
		return err
	}
	f, err := os.Create(path.Join(cwd, appConfig.AppName, "application.yml"))
	if err != nil {
		return err
	}
	if err = t.Execute(f, appConfig); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}
