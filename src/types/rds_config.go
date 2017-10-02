package types

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/pkg/errors"
)

// RdsConfig holds configuration fields for an rds dependency
type RdsConfig struct {
	AllocatedStorage   string             `yaml:"allocated-storage,omitempty"`
	InstanceClass      string             `yaml:"instance-class,omitempty"`
	DbName             string             `yaml:"db-name,omitempty"`
	Username           string             `yaml:",omitempty"`
	PasswordSecretName string             `yaml:"password-secret-name,omitempty"` //TODO: verify this is set via exo configure
	StorageType        string             `yaml:"storage-type,omitempty"`
	ServiceEnvVarNames ServiceEnvVarNames `yaml:"service-env-var-names,omitempty"`
}

// ServiceEnvVarNames are the names of RDS related env vars injected into a service
type ServiceEnvVarNames struct {
	DbName   string `yaml:"db-name,omitempty"`
	Username string `yaml:",omitempty"`
	Password string `yaml:",omitempty"`
}

// IsEmpty returns true if the given config object is empty
func (d *RdsConfig) IsEmpty() bool {
	return d.AllocatedStorage == "" &&
		d.InstanceClass == "" &&
		d.DbName == "" &&
		d.Username == "" &&
		d.PasswordSecretName == "" &&
		d.StorageType == "" &&
		d.ServiceEnvVarNames.DbName == "" &&
		d.ServiceEnvVarNames.Username == "" &&
		d.ServiceEnvVarNames.Password == ""
}

// ValidateFields validates that an rds config contains all required fields
func (d *RdsConfig) ValidateFields() error {
	requiredFields := []string{"AllocatedStorage", "InstanceClass", "DbName", "Username", "PasswordSecretName", "StorageType", "ServiceEnvVarNames"}
	missingField := checkRequiredFields(*d, requiredFields)
	if missingField != "" {
		return fmt.Errorf("missing required field 'Rds.%s'", missingField)
	}
	dbNameRegex := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	if !dbNameRegex.MatchString(d.DbName) {
		return errors.New("only alphanumeric characters and hyphens allowed in rds.db-name")
	}
	requiredServiceEnvVars := []string{"DbName", "Username", "Password"}
	missingField = checkRequiredFields(d.ServiceEnvVarNames, requiredServiceEnvVars)
	if missingField != "" {
		return fmt.Errorf("missing required field 'Rds.ServiceEnvVarNames.%s", missingField) //TODO address camelcase issue
	}
	return nil
}

func checkRequiredFields(object interface{}, requiredFields []string) string {
	for _, field := range requiredFields {
		value := reflect.ValueOf(object).FieldByName(field).String()
		if value == "" {
			return field
		}
	}
	return ""
}
