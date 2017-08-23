package compose

// ImagesOptions is the options to compose functions that deal with multiple images
type ImagesOptions struct {
	DockerComposeDir string
	Env              []string
	LogChannel       chan string
	ImageNames       []string
}
