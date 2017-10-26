package compose

import (
	"io"
)

// ImagesOptions is the options to compose functions that deal with multiple images
type ImagesOptions struct {
	DockerComposeDir string
	Env              []string
	ImageNames       []string
	Writer           io.Writer
	AbortOnExit      bool
}
