package dockerCompose

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// RunAllImages runs all the docker images
func RunAllImages(env []string, cwd string, write func(string)) (*exec.Cmd, *bytes.Buffer, error) {
	// cwd, , write, "Failed to run images"
	cmd, _, stdoutBuffer, err := start(cwd, env, "docker-compose", "up")
	return cmd, stdoutBuffer, err

}

// KillAllContainers kills all the containers
func KillAllContainers(env []string, cwd string, write func(string)) (*exec.Cmd, *bytes.Buffer, error) {
	// start([]string{"docker-compose", "down"}, cwd, env, write, "Failed to kill containers")
	cmd, _, stdoutBuffer, err := start(cwd, env, "docker-compose", "down")
	return cmd, stdoutBuffer, err
}

// func start(command []string, cwd string, env []string, write func(string), errorMessage string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
// 	fmt.Println(command, cwd)
// 	cmd := exec.Command(command[0], command[1:]...) // nolint gas
// 	cmd.Env = env
// 	cmd.Dir = cwd
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		log.Fatalf("Failed to set up a pipe for stdout: %s", err)
// 	}
// 	if err := cmd.Start(); err != nil {
// 		log.Fatalf("%s: %s", errorMessage, err)
// 	}
// 	in := bufio.NewScanner(stdout)
// 	for in.Scan() {
// 		write(in.Text())
// 	}
// }

// Start runs the given command in the given dir directory, and returns
// the pointer to the command, stdout pipe, output buffer and error (if any)
func start(dir string, env []string, command ...string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	cmd := exec.Command(command[0], command[1:]...) // nolint gas
	cmd.Dir = dir
	if len(env) > 0 {
		cmd.Env = env
	}
	in, err := cmd.StdinPipe()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err != nil {
		return nil, in, &out, err
	}
	if err = cmd.Start(); err != nil {
		return nil, in, &out, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return cmd, in, &out, nil
}
