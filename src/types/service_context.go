package types

// ServiceContext represents the exosphere service the user is running
type ServiceContext struct {
	Dir         string
	Location    string
	Config      ServiceConfig
	AppContext  AppContext
	ServiceData *ServiceData
}
