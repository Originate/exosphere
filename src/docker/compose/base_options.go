package compose

import (
	"io"
)

// BaseOptions are the options passed into docker compose functions
type BaseOptions struct {
	DockerComposeDir      string
	DockerComposeFileName string
	Env                   []string
	Writer                io.Writer
}
