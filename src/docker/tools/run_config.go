package tools

import "io"

//RunConfig holds fields for possible configuration passed to a 'docker run' command
type RunConfig struct {
	Volumes     []string
	Interactive bool
	WorkingDir  string
	ImageName   string
	Command     []string
	Writer      io.Writer
}
