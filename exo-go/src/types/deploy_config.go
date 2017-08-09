package types

type DeployConfig struct {
	AppConfig      AppConfig
	ServiceConfigs map[string]ServiceConfig
	Logger         chan string
	AppDir         string
	HomeDir        string
	AwsConfig      AwsConfig
}
