package cmd

import (
	"os"

	"github.com/Originate/exosphere/src/types/context"
)

// GetUserContext returns a UserContext for the current working direcotry
func GetUserContext() (*context.UserContext, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	appContext, err := context.GetAppContext(currentDir)
	if err != nil {
		return nil, err
	}
	if _, err = os.Stat("service.yml"); err != nil {
		return &context.UserContext{AppContext: appContext}, nil
	}
	serviceContext, err := appContext.GetInternalServiceContext(currentDir)
	if err != nil {
		return nil, err
	}
	return &context.UserContext{
		AppContext:        appContext,
		ServiceContext:    *serviceContext,
		HasServiceContext: true,
	}, nil
}
