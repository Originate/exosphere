package types

// ProductionDependencyConfigOptions represents the configuration of a development dependency
type ProductionDependencyConfigOptions struct {
	Rds RdsConfig `yaml:",omitempty"`
}
