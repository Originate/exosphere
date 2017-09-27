package types

import (
	"fmt"
	"os"
	"path"
)

type AppContext struct {
	Name     string
	Location string
	Config   AppConfig
}

func GetAppContext(location string) (AppContext, error) {
	walkDir := path.Clean(location)
	if !path.IsAbs(location) {
		return AppContext{}, fmt.Errorf("Could not get app context: location must be absolute")
	}
	for {
		if _, err := os.Stat(path.Join(walkDir, "application.yml")); err == nil {
			appConfig, err := NewAppConfig(walkDir)
			if err != nil {
				return AppContext{}, err
			}
			return AppContext{
				Name:     path.Base(walkDir),
				Location: walkDir,
				Config:   appConfig,
			}, nil
		}
		if walkDir == "/" || walkDir == "." {
			return AppContext{}, fmt.Errorf("%s is not an exosphere application", location)
		}
		walkDir = path.Dir(walkDir)
	}
}
