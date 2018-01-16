package aws

// Options is the structure for options to most functions
type Options struct {
	BucketName         string
	Profile            string
	Region             string
	TerraformLockTable string
}
