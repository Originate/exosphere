package types

// ProductionDependencyConfigOptions represents the configuration of a development dependency
type ProductionDependencyConfigOptions struct {
	AllocatedStorage string `yaml:"allocated-storage,omitempty"`
	InstanceClass    string `yaml:"instance-class,omitempty"`
	DbName           string `yaml:"db-name,omitempty"`
	Username         string `yaml:",omitempty"`
	PasswordEnvVar   string `yaml:"password-env-var,omitempty"` //TODO: verify this is set via exo configure
	StorageType      string `yaml:"storage-type,omitempty"`
}

// IsEmpty returns true if the given dependencyConfig object is empty
func (d *ProductionDependencyConfigOptions) IsEmpty() bool {
	return true
}
