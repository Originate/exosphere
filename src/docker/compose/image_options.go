package compose

// ImageOptions is the options to compose functions that deal with a single image
type ImageOptions struct {
	DockerComposeDir string
	Env              []string
	LogChannel       chan string
	ImageName        string
}
