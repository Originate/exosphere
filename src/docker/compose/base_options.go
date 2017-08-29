package compose

// BaseOptions are the options passed into docker compose functions
type BaseOptions struct {
	DockerComposeDir string
	Env              []string
	LogChannel       chan string
}
