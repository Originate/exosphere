package types

// AwsConfig contains top level information about an application's AWS account
type AwsConfig struct {
	Region               string
	AccountID            string
	SslCertificateArn    string
	Profile              string
	SecretsBucket        string
	TerraformStateBucket string
	TerraformLockTable   string
}
