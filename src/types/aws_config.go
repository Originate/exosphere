package types

// AwsConfig contains top level information about an application's AWS account
type AwsConfig struct {
	AccountID          string
	BucketName         string
	Profile            string
	Region             string
	SslCertificateArn  string
	TerraformLockTable string
}
