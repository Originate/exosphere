package config

// TerraformConfig contains information needed to
// generate and execute Terraform scripts
type TerraformConfig struct {
	AppConfig      AppConfig
	ServiceConfigs map[string]ServiceConfig
	AppDir         string
	HomeDir        string
	TerraformDir   string
	RemoteBucket   string
	LockTable      string
	Region         string
}
