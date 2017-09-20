package types

// ProductionDependencyConfigOptions represents the configuration of a development dependency
type ProductionDependencyConfigOptions struct {
	//TODO: populate with RDS fields
}

// IsEmpty returns true if the given dependencyConfig object is empty
func (d *ProductionDependencyConfigOptions) IsEmpty() bool {
	return true
}
