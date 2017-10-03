package types

import "os"

// Context represents contextual information about what application/service the user is running
type Context struct {
	AppContext        AppContext
	ServiceContext    ServiceContext
	HasServiceContext bool
}

// GetContext returns a Context for the current working direcotry
func GetContext() (Context, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return Context{}, err
	}
	appContext, err := GetAppContext(currentDir)
	if err != nil {
		return Context{}, err
	}
	if _, err = os.Stat("service.yml"); err != nil {
		return Context{AppContext: appContext}, nil
	}
	serviceContext, err := GetServiceContext(appContext, currentDir)
	if err != nil {
		return Context{}, err
	}
	return Context{
		AppContext:        appContext,
		ServiceContext:    serviceContext,
		HasServiceContext: true,
	}, nil
}
