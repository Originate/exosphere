package testHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	shellwords "github.com/mattn/go-shellwords"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

func emptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, os.FileMode(0777))
}

func run(command string) (string, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	outputArray, err := cmd.CombinedOutput()
	output := string(outputArray)
	return output, err
}

func start(command string, dir string) (io.WriteCloser, bytes.Buffer, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return nil, bytes.Buffer{}, err
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	cmd.Dir = dir
	stdinPipe, err := cmd.StdinPipe()
	var stdoutBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	if err != nil {
		return stdinPipe, stdoutBuffer, err
	}
	if err = cmd.Start(); err != nil {
		return stdinPipe, stdoutBuffer, fmt.Errorf("Error running %s\nError:%s", command, err)
	}
	return stdinPipe, stdoutBuffer, nil
}

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}

func waitForText(stdout bytes.Buffer, text string, duration int) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(time.Duration(duration) * time.Millisecond)
	for !strings.Contains(stdout.String(), text) {
		select {
		case <-ticker.C:
			return nil
		case <-timeout:
			return fmt.Errorf("Timed out after %d milliseconds", duration)
		}
	}
	return nil
}
