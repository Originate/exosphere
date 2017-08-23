package compose

// BaseOptions is the options to compose functions that deal with all images
type BaseOptions struct {
	DockerComposeDir string
	Env              []string
	LogChannel       chan string
}
