package types

// ProductionDependencyConfigOptions represents the configuration of a development dependency
type ProductionDependencyConfigOptions struct {
	Rds RdsConfig `yaml:",omitempty"`
}

// IsEmpty returns true if the given config object is empty
func (d *ProductionDependencyConfigOptions) IsEmpty() bool {
	return d.Rds.IsEmpty()
}
