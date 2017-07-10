package dockerCompose

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

// RunAllImages runs all the docker images
func RunAllImages(env []string, cwd string, write func(string)) (*bytes.Buffer, error) {
	return start(cwd, env, write, "docker-compose", "up")
}

// KillAllContainers kills all the containers
func KillAllContainers(env []string, cwd string, write func(string)) (*bytes.Buffer, error) {
	return start(cwd, env, write, "docker-compose", "down")
}

func start(dir string, env []string, write func(string), command ...string) (*bytes.Buffer, error) {
	cmd := exec.Command(command[0], command[1:]...) // nolint gas
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = env
	}
	var out bytes.Buffer
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &out, err
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			write(text)
			out.Write([]byte(text))
		}
	}()
	if err = cmd.Start(); err != nil {
		return &out, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return &out, nil
}
