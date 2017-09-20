package types

// DeployConfig contains information needed for deployment
type DeployConfig struct {
	AppConfig                AppConfig
	ServiceConfigs           map[string]ServiceConfig
	LogChannel               chan string
	AppDir                   string
	DockerComposeProjectName string
	TerraformDir             string
	SecretsPath              string
	AwsConfig                AwsConfig
}
