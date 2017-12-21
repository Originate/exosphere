package types

// AppRemoteConfig represents production specific configuration for an application
type AppRemoteConfig struct {
	Dependencies map[string]RemoteDependency
	Environments map[string]AppRemoteEnvironment
}

// ValidateFields validates that this section contains the required fields
func (p AppRemoteConfig) ValidateFields(remoteEnvironmentID string) error {
	return p.Environments[remoteEnvironmentID].ValidateFields(remoteEnvironmentID)
}
