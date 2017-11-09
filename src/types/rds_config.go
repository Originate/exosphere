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

// yamlNames maps struct field names to their names in a yaml file
// used in error messages given to the user
var yamlNames = map[string]string{
	"AllocatedStorage":   "allocated-storage",
	"InstanceClass":      "instance-class",
	"DbName":             "db-name",
	"Username":           "username",
	"Password":           "password",
	"PasswordSecretName": "password-secret-name",
	"StorageType":        "storage-type",
	"ServiceEnvVarNames": "service-env-var-names",
}

// ValidateFields validates that an rds config contains all required fields
func (d *RdsConfig) ValidateFields() error {
	requiredFields := []string{"AllocatedStorage", "InstanceClass", "DbName", "Username", "PasswordSecretName", "StorageType", "ServiceEnvVarNames"}
	missingField := checkRequiredFields(*d, requiredFields)
	if missingField != "" {
		return fmt.Errorf("missing required field 'rds.%s'", yamlNames[missingField])
	}
	dbNameRegex := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	if !dbNameRegex.MatchString(d.DbName) {
		return errors.New("only alphanumeric characters and hyphens allowed in 'rds.db-name'")
	}
	requiredServiceEnvVars := []string{"DbName", "Username", "Password"}
	missingField = checkRequiredFields(d.ServiceEnvVarNames, requiredServiceEnvVars)
	if missingField != "" {
		return fmt.Errorf("missing required field 'rds.service-env-var-names.%s", yamlNames[missingField])
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
