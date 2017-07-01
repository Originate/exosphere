package testHelpers

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

const validateTextContainsErrorTemplate = `
Expected:

%s

to include

%s
	`

func EmptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, 0777)
}

func ReformatCommand(command string) []string {
	return strings.Split(command, " ")
}

func ValidateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}

func WaitForText(stdout bytes.Buffer, text string, duration int) error {
	interval := time.Tick(100 * time.Millisecond)
	timeout := time.After(time.Duration(duration) * time.Millisecond)
	for !strings.Contains(stdout.String(), text) {
		select {
		case <-interval:
			return nil
		case <-timeout:
			return fmt.Errorf("Timed out after %s milliseconds", duration)
		}
	}
	return nil
}
