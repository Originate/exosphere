package types

// ServiceDependencyData represents data about a service to be passed to a dependencies
// The first key is the dependency name, the value is data to the passed to the dependency
type ServiceDependencyData map[string]map[string]interface{}
