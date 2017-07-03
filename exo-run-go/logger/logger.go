package logger

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/exo-run-go/helpers"
	"github.com/fatih/color"
)

type Logger struct {
	Roles         []string
	SilencedRoles []string
	Length        int
}

func (logger *Logger) Colors() string {
	return map[string]interface{}{"exocom": color.Cyan}
}

func (logger *Logger) Log(role, text string, trim bool) {
	if trim {
		text := strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !helpers.Contains(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", color.Bold(logger.Pad(left)), right)
			if color, exists := logger.Colors()[left]; exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
}

func (logger *Logger) Error(role, text string, trim bool) {
	if trim {
		text := strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !helpers.Contains(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", color.Bold(logger.Pad(left)), right)
			if color, exists := logger.Colors[left]; exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
}

func (logger *Logger) SetColors(roles) {
	defaultColors := []interface{}{color.Magenta, color.Blue, color.Yellow, color.Cyan}
	for i, role := range roles {
		Logger.Colors[role] = defaultColors[i%len(defaultColors)]
		Logger.Length = math.Max(Logger.Length, len(role))
	}
}

func (logger *Logger) Pad(text string) {
	return padLeft(text, logger.Length, ' ')
}

func parseLine(role, line string) (string, string) {
	segments := []string{}
	for _, segment := range regexp.MustCompile(`\s+\|\s*`).Split(line, -1) {
		segments = append(segments, strings.TrimSpace(segment))
	}
	if len(segments) == 2 {
		service := parseService(segments[0])
		return service, reformatLine(segments[1])
	}
	return role, line
}

func parseService(text string) {
	return stripColor(strip(`(\d+\.)?(\d+\.)?(\*|\d+)$`, text))
}

func reformatLine(line string) {
	return strings.TrimSpace(stripColor(line))
}

func padLeft(text string, length int, padding rune) {
	return fmt.Sprintf("%s%s", strings.Repeat(padding, length-len(text)), text)
}

func stripColor(text string) {
	return strip(`\033\[[0-9;]*m`, text)
}

func strip(regex, text string) {
	return regexp.MustCompile(regex).ReplaceAll(text, "")
}
