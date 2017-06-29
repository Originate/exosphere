package testHelpers

import (
	"fmt"
	"os"
	"path"
	"strings"
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
	return os.Mkdir(dir, os.FileMode(0777))
}

func ReformatCommand(command, cwd string) []string {
	return strings.Split(path.Join(cwd, "bin", command), " ")
}

func ValidateTextContains(haystack, needle string) error {
	if strings.Contains(haystack, needle) {
		return nil
	}
	return fmt.Errorf(validateTextContainsErrorTemplate, haystack, needle)
}
