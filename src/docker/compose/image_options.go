package compose

import "github.com/Originate/exosphere/src/util"

// ImageOptions is the options to compose functions that deal with a single image
type ImageOptions struct {
	DockerComposeDir string
	Env              []string
	ImageName        string
	Logger           *util.Logger
}
