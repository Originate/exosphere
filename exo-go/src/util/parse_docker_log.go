package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/stringplus"
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

// ParseDockerComposeLog parses the given docker-compose output
// and returns the service name and service output
func ParseDockerComposeLog(role, line string) (string, string) {
	segments := []string{}
	for _, segment := range regexp.MustCompile(`\s+\|\s*`).Split(line, -1) {
		segments = append(segments, strings.TrimSpace(segment))
	}
	if len(segments) == 2 {
		return parseService(segments[0]), reformatLine(segments[1])
	}
	return role, line
}

func parseService(text string) string {
	return stringplus.StripAnsiColors(stringplus.StripRegexp(`(\d+\.)?(\d+\.)?(\*|\d+)$`, text))
}

func reformatLine(line string) string {
	return strings.TrimSpace(stringplus.StripAnsiColors(line))
}
