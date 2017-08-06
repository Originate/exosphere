package docker

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/runplus"
)

// RunInDockerImage runs the given command in a new writeable container layer
// over the given image, removes the container when the command exits, and returns
// the output string and an error if any
func RunInDockerImage(imageName, command string) (string, error) {
	return runplus.Run("", fmt.Sprintf("docker run --rm %s %s", imageName, command))
}
