package composebuilder

// BuildMode is what type of docker compose config should be created
type BuildMode uint

const (
	// BuildModeLocalDevelopment used for `exo run` without production flag
	BuildModeLocalDevelopment BuildMode = iota
	// BuildModeLocalProduction used for `exo run` with production flag
	BuildModeLocalProduction
	// BuildModeDeployProduction used for `exo deploy`
	BuildModeDeployProduction
)
