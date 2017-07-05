package dockerCompose

import (
	"bufio"
	"log"
	"os/exec"
)

func RunAllImages(env map[string]interface{}, cwd string, write func(string)) {
	run([]string{"docker-compose", "up"}, cwd, write, "Failed to run images")
}

func KillAllContainers(env map[string]interface{}, cwd string, write func(string)) {
	run([]string{"docker-compose", "down"}, cwd, write, "Failed to kill containers")
}

func run(command []string, dir string, write func(string), errorMessage string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to set up a pipe for stdout: %s", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("%s: %s", errorMessage, err)
	}
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		write(in.Text())
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("%s: %s", errorMessage, err)
	}
}
