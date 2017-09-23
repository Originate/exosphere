package composebuilder

// BuildMode is what type of docker compose config should be created
type BuildMode uint

const (
	// BuildModeLocalDevelopment used for `exo run` and `exo test` without flags
	BuildModeLocalDevelopment BuildMode = iota
	// BuildModeLocalDevelopmentNoMount used for `exo run` and `exo test` with no-mount flag
	BuildModeLocalDevelopmentNoMount
	// BuildModeLocalProduction used for `exo run` with production flag
	BuildModeLocalProduction
	// BuildModeDeployProduction used for `exo deploy`
	BuildModeDeployProduction
)
