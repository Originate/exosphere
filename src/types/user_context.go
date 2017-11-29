package types

import (
	"fmt"
	"os"
	"path"
)

// UserContext represents contextual information about what application/service the user is running
type UserContext struct {
	AppContext        *AppContext
	ServiceContext    *ServiceContext
	HasServiceContext bool
}

// GetAppContext returns an AppContext for the application found at the passed location
// by walkig up the directory tree until it finds an application.yml
func GetAppContext(location string) (*AppContext, error) {
	walkDir := path.Clean(location)
	if !path.IsAbs(location) {
		return nil, fmt.Errorf("Could not get app context: location must be absolute")
	}
	for {
		if _, err := os.Stat(path.Join(walkDir, "application.yml")); err == nil {
			appConfig, err := NewAppConfig(walkDir)
			if err != nil {
				return nil, err
			}
			return &AppContext{
				Location: walkDir,
				Config:   appConfig,
			}, nil
		}
		if walkDir == "/" || walkDir == "." {
			return nil, fmt.Errorf("%s is not an exosphere application", location)
		}
		walkDir = path.Dir(walkDir)
	}
}
