package types

// AwsConfig contains top level information about an application's AWS account
type AwsConfig struct {
	Region               string
	Profile              string
	CredentialsFile      string
	SecretsBucket        string
	TerraformStateBucket string
	TerraformLockTable   string
}
