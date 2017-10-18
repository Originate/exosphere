package cmd

import (
	"os"

	"github.com/Originate/exosphere/src/types"
)

// GetContext returns a Context for the current working direcotry
func GetContext() (types.Context, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return types.Context{}, err
	}
	appContext, err := types.GetAppContext(currentDir)
	if err != nil {
		return types.Context{}, err
	}
	if _, err = os.Stat("service.yml"); err != nil {
		return types.Context{AppContext: appContext}, nil
	}
	serviceContext, err := appContext.GetServiceContext(currentDir)
	if err != nil {
		return types.Context{}, err
	}
	return types.Context{
		AppContext:        appContext,
		ServiceContext:    serviceContext,
		HasServiceContext: true,
	}, nil
}
