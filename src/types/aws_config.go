package types

// AwsConfig contains top level information about an application's AWS account
type AwsConfig struct {
	Region               string
	SecretsBucket        string
	TerraformStateBucket string
	TerraformLockTable   string
}
