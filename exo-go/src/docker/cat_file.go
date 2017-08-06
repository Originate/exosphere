package docker

import (
	"fmt"

	"github.com/moby/moby/client"
)

// CatFileInDockerImage reads the file fileName inside the given image
func CatFileInDockerImage(c *client.Client, imageName, fileName string) ([]byte, error) {
	if err := PullImage(c, imageName); err != nil {
		return []byte(""), err
	}
	command := fmt.Sprintf("cat %s", fileName)
	output, err := RunInDockerImage(imageName, command)
	return []byte(output), err
}
