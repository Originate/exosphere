package compose

import (
	"io"
)

// CommandOptions is the options to compose functions that deal with multiple images
type CommandOptions struct {
	DockerComposeDir      string
	DockerComposeFileName string
	Env                   []string
	ImageNames            []string
	Writer                io.Writer
	AbortOnExit           bool
	Build                 bool
}
