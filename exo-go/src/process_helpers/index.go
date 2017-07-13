package processHelpers

import (
	"os/exec"
)

// Run runs the given command, waits for the process to finish and
// returns the output string and error (if any)
func Run(dir string, commandWords ...string) (string, error) {
	cmd := exec.Command(commandWords[0], commandWords[1:]...) // nolint gass
	cmd.Dir = dir
	outputArray, err := cmd.CombinedOutput()
	output := string(outputArray)
	return output, err
}
