package types

// AppConfig represents the configuration of an application
type AppConfig struct {
	Name         string
	Description  string
	Version      string
	Dependencies []Dependency
	Services
	Templates map[string]string `yaml:",omitempty"`
}
