package types

// Dependency is an unexpected type
type Dependency struct {
	Name    string
	Version string
	//Config  DependencyConfig
}

// DependencyConfig is an unexpected type
type DependencyConfig struct {
	Ports                 []string
	Volumes               []string
	OnlineText            string
	DependencyEnvironment map[string]string
	ServiceEnvironment    map[string]string
}

// Service represents a service
type Service struct {
	Location string
}

// ServiceConfig represents the configuration of a service
type ServiceConfig struct {
	ServiceRole, ServiceType, Description, Author, TemplatePath, ModelName, ProtectionLevel string
}

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services     map[string]map[string]Service
}
