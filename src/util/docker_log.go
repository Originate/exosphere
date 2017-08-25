package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	matches := regexp.MustCompile(`^(\S*)\s+\|(.*)$`).FindStringSubmatch(line)
	if len(matches) == 3 {
		return parseService(matches[1]), strings.TrimPrefix(matches[2], " ")
	}
	return role, line
}

// NormalizeDockerComposeLog removes colors, control characters, and converts
// carriage returns to newlines
func NormalizeDockerComposeLog(line string) string {
	return strings.TrimSpace(convertCarriageReturnToNewline(stripControl(stripColor(line))))
}

func parseService(text string) string {
	return Strip(`(\d+\.)?(\d+\.)?(\*|\d+)$`, text)
}

// Strip strips off substrings that match the given regex from the
// given text
func Strip(regex, text string) string {
	return string(regexp.MustCompile(regex).ReplaceAll([]byte(text), []byte("")))
}

func stripColor(text string) string {
	return Strip(`\033\[[0-9;]*m`, text)
}

func stripControl(text string) string {
	return Strip(`\033\[(1A|2K|1B)`, text)
}

func convertCarriageReturnToNewline(text string) string {
	return string(regexp.MustCompile("\r").ReplaceAll([]byte(text), []byte("\n")))
}
