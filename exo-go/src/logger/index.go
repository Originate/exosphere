package logger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/fatih/color"
)

// Logger represents a logger
type Logger struct {
	Roles         []string
	SilencedRoles []string
	Length        int
	Colors        map[string]func(string, ...interface{})
}

// NewLogger is Logger's constructor
func NewLogger(roles, silencedRoles []string) *Logger {
	logger := &Logger{Roles: roles, SilencedRoles: silencedRoles}
	logger.Colors = map[string]func(string, ...interface{}){
		"exocom":     color.Cyan,
		"exo-run":    color.Green,
		"exo-clone":  color.Green,
		"exo-setup":  color.Green,
		"exo-test":   color.Green,
		"exo-sync":   color.Green,
		"exo-lint":   color.Green,
		"exo-deploy": color.Green,
	}
	logger.setColors(roles)
	return logger
}

func (logger *Logger) getColor(role string) (func(string, ...interface{}), bool) {
	chosenColor, exists := logger.Colors[role]
	return chosenColor, exists
}

func (logger *Logger) setColors(roles []string) {
	defaultColors := []func(string, ...interface{}){color.Magenta, color.Blue, color.Yellow, color.Cyan}
	for i, role := range roles {
		logger.Colors[role] = defaultColors[i%len(defaultColors)]
	}
	for role := range logger.Colors {
		if len(role) > logger.Length {
			logger.Length = len(role)
		}
	}
}

func (logger *Logger) pad(text string) string {
	return padLeft(text, logger.Length, ' ')
}

// Log logs the given text
func (logger *Logger) Log(role, text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !util.DoesStringArrayContain(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", logger.pad(left), right)
			if color, exists := logger.getColor(left); exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
}

// Error logs the given error text
func (logger *Logger) Error(role, text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !util.DoesStringArrayContain(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", logger.pad(left), right)
			if color, exists := logger.getColor(left); exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
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

func parseService(text string) string {
	return stripColor(strip(`(\d+\.)?(\d+\.)?(\*|\d+)$`, text))
}

func reformatLine(line string) string {
	return strings.TrimSpace(stripColor(line))
}

func padLeft(text string, length int, padding rune) string {
	return fmt.Sprintf("%s%s", strings.Repeat(string(padding), length-len(text)), text)
}

func stripColor(text string) string {
	return strip(`\033\[[0-9;]*m`, text)
}

func strip(regex, text string) string {
	return string(regexp.MustCompile(regex).ReplaceAll([]byte(text), []byte("")))
}
