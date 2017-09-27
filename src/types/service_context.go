package types

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	Name       string
	Location   string
	Config     ServiceConfig
	AppContext AppContext

	ServiceData *ServiceData
}

// GetServiceContext returns a ServiceContext for the service found at the passed location
func GetServiceContext(appContext AppContext, location string) (ServiceContext, error) {
	serviceConfig := ServiceConfig{}
	yamlFile, err := ioutil.ReadFile(path.Join(location, "service.yml"))
	if err != nil {
		return ServiceContext{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return ServiceContext{}, err
	}
	return ServiceContext{
		Name:       serviceConfig.Type,
		Location:   location,
		Config:     serviceConfig,
		AppContext: appContext,
	}, nil
}
