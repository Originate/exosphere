package cmd

import (
	"os"

	"github.com/Originate/exosphere/src/types"
)

// GetUserContext returns a UserContext for the current working direcotry
func GetUserContext() (*types.UserContext, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	appContext, err := types.GetAppContext(currentDir)
	if err != nil {
		return nil, err
	}
	if _, err = os.Stat("service.yml"); err != nil {
		return &types.UserContext{AppContext: appContext}, nil
	}
	serviceContext, err := appContext.GetServiceContext(currentDir)
	if err != nil {
		return nil, err
	}
	return &types.UserContext{
		AppContext:        appContext,
		ServiceContext:    serviceContext,
		HasServiceContext: true,
	}, nil
}
