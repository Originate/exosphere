package util

import (
	"regexp"
	"strings"
)

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
	return stripColor(strip(`(\d+\.)?(\d+\.)?(\*|\d+)$`, text))
}

func reformatLine(line string) string {
	return strings.TrimSpace(stripColor(line))
}

func strip(regex, text string) string {
	return string(regexp.MustCompile(regex).ReplaceAll([]byte(text), []byte("")))
}

func stripColor(text string) string {
	return strip(`\033\[[0-9;]*m`, text)
}
