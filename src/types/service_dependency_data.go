package types

import "github.com/Originate/exosphere/src/util"

// ServiceDependencyData represents data about a service to be passed to a dependencies
// The first key is the dependency name, the value is data to the passed to the dependency
type ServiceDependencyData map[string]map[string]interface{}

// StringifyMapKeys converts any maps in the data to having keys of string instead of interface
func (s ServiceDependencyData) StringifyMapKeys() {
	for _, v1 := range s {
		for k2, v2 := range v1 {
			v1[k2] = util.StringifyKeysMapValue(v2)
		}
	}
}
