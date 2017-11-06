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

// LocalDevelopmentComposeFileName is the docker-compose file name for local development runs
const LocalDevelopmentComposeFileName = "run_development.yml"

// LocalProductionComposeFileName is the docker-compose file name for local production runs
const LocalProductionComposeFileName = "run_production.yml"

// LocalTestComposeFileName is the docker-compose file name for local test runs
const LocalTestComposeFileName = "test.yml"

// GetDockerComposeFileName returns the proper docker-compose file name for the build environment
func (b BuildMode) GetDockerComposeFileName() string {
	if b.Type == BuildModeTypeLocal {
		switch b.Environment {
		case BuildModeEnvironmentDevelopment:
			return LocalDevelopmentComposeFileName
		case BuildModeEnvironmentProduction:
			return LocalProductionComposeFileName
		case BuildModeEnvironmentTest:
			return LocalTestComposeFileName
		}
	}
	return "docker-compose.yml"
}

// GetLocalRunComposeFileNames returns a list of docker-compose file names for local run processes
func GetLocalRunComposeFileNames() []string {
	return []string{LocalDevelopmentComposeFileName, LocalProductionComposeFileName}
}
