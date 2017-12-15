package types

import (
	"fmt"
	"reflect"
)

// AppRemoteConfigEnvironment represents production specific configuration for an application
type AppRemoteConfigEnvironment struct {
	URL               string `yaml:",omitempty"`
	Region            string `yaml:",omitempty"`
	AccountID         string `yaml:"account-id,omitempty"`
	SslCertificateArn string `yaml:"ssl-certificate-arn,omitempty"`
}

// ValidateFields validates that the production section contiains the required fields
func (a AppRemoteConfigEnvironment) ValidateFields(remoteEnvironmentID string) error {
	requiredFields := []string{"URL", "Region", "AccountID", "SslCertificateArn"}
	for _, field := range requiredFields {
		value := reflect.ValueOf(a).FieldByName(field).String()
		if value == "" {
			return fmt.Errorf("application.yml missing required field 'remote.environments.%s.%s'", remoteEnvironmentID, field)
		}
	}
	return nil
}
