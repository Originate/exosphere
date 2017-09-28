package compose

import "github.com/Originate/exosphere/src/util"

// BaseOptions are the options passed into docker compose functions
type BaseOptions struct {
	DockerComposeDir string
	Env              []string
	Logger           *util.Logger
}
