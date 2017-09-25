package types

// RdsConfig holds configuration fields for an rds dependency
type RdsConfig struct {
	AllocatedStorage string `yaml:"allocated-storage,omitempty"`
	InstanceClass    string `yaml:"instance-class,omitempty"`
	DbName           string `yaml:"db-name,omitempty"`
	Username         string `yaml:",omitempty"`
	PasswordEnvVar   string `yaml:"password-env-var,omitempty"` //TODO: verify this is set via exo configure
	StorageType      string `yaml:"storage-type,omitempty"`
}

// IsEmpty returns true if the given config object is empty
func (d *RdsConfig) IsEmpty() bool {
	return d.AllocatedStorage == "" &&
		d.InstanceClass == "" &&
		d.DbName == "" &&
		d.Username == "" &&
		d.PasswordEnvVar == "" &&
		d.StorageType == ""
}
