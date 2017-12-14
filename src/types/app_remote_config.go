package types

import (
	"fmt"
	"reflect"
)

// AppRemoteConfig represents production specific configuration for an application
type AppRemoteConfig struct {
	Dependencies      map[string]RemoteDependency
	Environment       map[string]string
	Secrets           []string
	URL               string `yaml:",omitempty"`
	Region            string `yaml:",omitempty"`
	AccountID         string `yaml:"account-id,omitempty"`
	SslCertificateArn string `yaml:"ssl-certificate-arn,omitempty"`
}

// ValidateFields validates that the production section contiains the required fields
func (p AppRemoteConfig) ValidateFields() error {
	requiredFields := []string{"URL", "Region", "AccountID", "SslCertificateArn"}
	for _, field := range requiredFields {
		value := reflect.ValueOf(p).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("application.yml missing required field 'remote.%s'", field)
		}
	}
	return nil
}
