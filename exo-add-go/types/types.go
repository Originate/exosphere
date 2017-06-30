package types

// Dependency is an unexported type
type Dependency struct {
	Name    string
	Version string
}

// Service represents a service
type Service struct {
	Location string
}

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services     map[string]map[string]Service
}
