package dockercompose

import (
	"fmt"
	"regexp"
	"strconv"
)

// GetServiceExitCode parses the given docker-compose output and
// returns the exit code of the given service and an error if any
func GetServiceExitCode(role, dockerComposeLog string) (int, error) {
	exitCodeRegex, err := regexp.Compile(fmt.Sprintf("%s exited with code (?P<first>[0-9])", role))
	if err != nil {
		return 1, err
	}
	matches := exitCodeRegex.FindStringSubmatch(dockerComposeLog)
	if len(matches) < 2 {
		return 1, fmt.Errorf("could not get exit code for %s", role)
	}
	return strconv.Atoi(matches[1])
}
