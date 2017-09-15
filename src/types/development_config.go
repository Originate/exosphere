package types

// DevelopmentConfig is the development key listed in a ServiceConfig
type DevelopmentConfig struct {
	Scripts map[string]string `yaml:",omitempty"`
}
