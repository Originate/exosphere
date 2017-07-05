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

// EmptyDir creates an empty dir directory
func EmptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Mkdir(dir, 0777)
}

// ValidateTextContains returns an error if needle is not in haystack, and
// nil otherwise
func ValidateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}

// WaitForText checks stdout every 100 milliseconds for the text, and returns
// an error if the wait time exceeds duration. It returns nil otherwise.
func WaitForText(stdout *bytes.Buffer, text string, duration int) error {
	interval, timeout := time.Tick(100*time.Millisecond), time.After(time.Duration(duration)*time.Millisecond)
	var output string
	for !strings.Contains(output, text) {
		select {
		case <-interval:
			output = stdout.String()
		case <-timeout:
			return fmt.Errorf("Timed out after %s milliseconds", duration)
		}
	}
	return nil
}
