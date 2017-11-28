package cmd

import (
	"os"

	"github.com/Originate/exosphere/src/types/context"
)

// GetUserContext returns a UserContext for the current working direcotry
func GetUserContext() (context.UserContext, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return context.UserContext{}, err
	}
	appContext, err := context.GetAppContext(currentDir)
	if err != nil {
		return context.UserContext{}, err
	}
	if _, err = os.Stat("service.yml"); err != nil {
		return context.UserContext{AppContext: appContext}, nil
	}
	serviceContext, err := appContext.GetServiceContext(currentDir)
	if err != nil {
		return context.UserContext{}, err
	}
	return context.UserContext{
		AppContext:        appContext,
		ServiceContext:    serviceContext,
		HasServiceContext: true,
	}, nil
}
