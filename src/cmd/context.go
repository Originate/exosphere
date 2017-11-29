package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Originate/exosphere/src/config"
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
	serviceContexts, err := config.GetServiceContexts(appContext)
	if err != nil {
		return nil, err
	}
	var matchingServiceContext *types.ServiceContext
	for _, serviceContext := range serviceContexts {
		if currentDir == path.Join(appContext.Location, serviceContext.AppData.Location) {
			matchingServiceContext = serviceContext
		}
	}
	if matchingServiceContext == nil {
		return nil, fmt.Errorf("Service is not listed in application.yml")
	}
	return &types.UserContext{
		AppContext:        appContext,
		ServiceContext:    matchingServiceContext,
		HasServiceContext: true,
	}, nil
}
