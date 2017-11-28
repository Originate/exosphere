package context

import "github.com/Originate/exosphere/src/types"

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	External    bool
	Dir         string
	Location    string
	Config      types.ServiceConfig
	AppContext  *AppContext
	ServiceData *types.ServiceData
}
