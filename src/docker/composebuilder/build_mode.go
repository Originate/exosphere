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

const LocalDevelopmentComposeFileName = "run_development.yml"
const LocalProductionComposeFileName = "run_production.yml"
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

// GetDevelopmentDockerComposeFileNames returns a list of local dev/prod compose file names
func GetDevelopmentDockerComposeFileNames() []string {
	return []string{LocalDevelopmentComposeFileName, LocalProductionComposeFileName}
}
