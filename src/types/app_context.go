package types

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// AppContext represents the exosphere application the user is running
type AppContext struct {
	Location string
	Config   AppConfig
}

// GetServiceContext returns a ServiceContext for the service found at the passed location
func (a AppContext) GetServiceContext(location string) (ServiceContext, error) {
	serviceConfig := ServiceConfig{}
	yamlFile, err := ioutil.ReadFile(path.Join(location, "service.yml"))
	if err != nil {
		return ServiceContext{}, err
	}
	if err = yaml.Unmarshal(yamlFile, &serviceConfig); err != nil {
		return ServiceContext{}, err
	}
	return ServiceContext{
		Dir:        path.Base(location),
		Location:   location,
		Config:     serviceConfig,
		AppContext: a,
	}, nil
}
