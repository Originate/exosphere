package dockerCompose

import (
	"fmt"
	"log"
	"time"

	"github.com/Originate/exosphere/exo-go/src/process_helpers"
)

// RunAllImages runs all the docker images
func RunAllImages(env []string, cwd string, write func(string)) {
	_, _, stdoutBuffer, err := processHelpers.Start("docker-compose up", cwd, env)
	if err != nil {
		log.Fatalf("Failed to run docker containers: %s", err)
	}
	write(getOutput(stdoutBuffer))
}

// KillAllContainers kills all existing docker containers
func KillAllContainers(env []string, cwd string, write func(string)) {
	_, _, stdoutBuffer, err := processHelpers.Start("docker-compose down", cwd, env)
	if err != nil {
		log.Fatalf("Failed to kill docker containers: %s", err)
	}
	write(getOutput(stdoutBuffer))
}

func getOutput(stdoutBuffer fmt.Stringer) string {
	timeout := time.After(time.Duration(1000) * time.Millisecond)
	select {
	case <-timeout:
		return stdoutBuffer.String()
	}
	return stdoutBuffer.String()
}
