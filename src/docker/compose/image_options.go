package compose

import (
	"io"
)

// ImageOptions is the options to compose functions that deal with a single image
type ImageOptions struct {
	DockerComposeDir string
	Env              []string
	ImageName        string
	Writer           io.Writer
}
