package testHelpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/DATA-DOG/godog/gherkin"
	shellwords "github.com/mattn/go-shellwords"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

func enterInput(in io.WriteCloser, out fmt.Stringer, row *gherkin.TableRow) error {
	field, input := row.Cells[0].Value, row.Cells[1].Value
	if err := waitForText(out, field, 1000); err != nil {
		return err
	}
	_, err := in.Write([]byte(input + "\n"))
	return err
}

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

func start(command string, dir string) (*exec.Cmd, io.WriteCloser, *bytes.Buffer, error) {
	commandWords, err := shellwords.Parse(command)
	if err != nil {
		return nil, nil, &bytes.Buffer{}, err
	}
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gas
	cmd.Dir = dir
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

func validateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}

func waitForText(stdout fmt.Stringer, text string, duration int) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(time.Duration(duration) * time.Millisecond)
	var output string
	for !strings.Contains(output, text) {
		select {
		case <-ticker.C:
			output = stdout.String()
		case <-timeout:
			return fmt.Errorf("Timed out after %d milliseconds", duration)
		}
	}
	return nil
}
