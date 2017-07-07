package dockerCompose

import (
	"bufio"
	"log"
	"os/exec"
)

func RunAllImages(env []string, cwd string, write func(string)) {
	run([]string{"docker-compose", "up"}, cwd, env, write, "Failed to run images")
}

func KillAllContainers(env []string, cwd string, write func(string)) {
	run([]string{"docker-compose", "down"}, cwd, env, write, "Failed to kill containers")
}

func run(command []string, cwd string, env []string, write func(string), errorMessage string) {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = env
	cmd.Dir = cwd
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
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
	in := bufio.NewScanner(stderr)
	for in.Scan() {
		write(in.Text())
	}
}
