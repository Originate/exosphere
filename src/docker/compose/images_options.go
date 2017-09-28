package compose

import "github.com/Originate/exosphere/src/util"

// ImagesOptions is the options to compose functions that deal with multiple images
type ImagesOptions struct {
	DockerComposeDir string
	Env              []string
	ImageNames       []string
	Logger           *util.Logger
}
