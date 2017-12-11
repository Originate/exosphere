package types

// RemoteDependencyConfig represents the configuration of a development dependency
type RemoteDependencyConfig struct {
	Rds     RdsConfig `yaml:",omitempty"`
	Version string    `yaml:",omitempty"`
}
