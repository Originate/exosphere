package types

import (
	"fmt"
	"reflect"
)

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

// ValidateFields validates that an rds config contains all required fields
func (d *RdsConfig) ValidateFields() error {
	requiredFields := []string{"DbName", "Username", "PasswordEnvVar"}
	for _, field := range requiredFields {
		value := reflect.ValueOf(*d).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("missing required field 'Rds.%s'", field)
		}
	}
	return nil
}
