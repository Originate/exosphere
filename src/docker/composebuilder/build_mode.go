package composebuilder

// BuildMode determines what type of docker compose config should be created
type BuildMode struct {
	Type        BuildModeType
	Mount       bool
	Environment BuildModeEnvironment
}

// BuildModeType indicates whether the docker compose config should be created local use or for deployment
type BuildModeType uint

// Possible values for BuildModeType
const (
	BuildModeTypeLocal = iota
	BuildModeTypeDeploy
)

// BuildModeEnvironment indicates which environment to build the docker compose config for
type BuildModeEnvironment uint

// Possible values for BuildModeEnvironment
const (
	BuildModeEnvironmentTest = iota
	BuildModeEnvironmentDevelopment
	BuildModeEnvironmentProduction
)

// GetDockerComposeFileName returns the proper docker-compose file name for the build environment
func (b BuildMode) GetDockerComposeFileName() string {
	switch b.Environment {
	case BuildModeEnvironmentDevelopment:
		return "run_development.yml"
	case BuildModeEnvironmentProduction:
		return "run_production.yml"
	case BuildModeEnvironmentTest:
		return "test.yml"
	default:
		return "docker-compose.yml"
	}
}
