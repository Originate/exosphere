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
	Colors        map[string]color.Attribute
}

// NewLogger is Logger's constructor
func NewLogger(roles, silencedRoles []string) *Logger {
	logger := &Logger{Roles: roles, SilencedRoles: silencedRoles}
	logger.Colors = map[string]color.Attribute{
		"exocom":     color.FgCyan,
		"exo-run":    color.FgGreen,
		"exo-clone":  color.FgGreen,
		"exo-setup":  color.FgGreen,
		"exo-test":   color.FgGreen,
		"exo-sync":   color.FgGreen,
		"exo-lint":   color.FgGreen,
		"exo-deploy": color.FgGreen,
	}
	logger.setColors(roles)
	return logger
}

func (logger *Logger) getColor(role string) (color.Attribute, bool) {
	chosenColor, exists := logger.Colors[role]
	return chosenColor, exists
}

func (logger *Logger) setColors(roles []string) {
	defaultColors := []color.Attribute{color.FgMagenta, color.FgBlue, color.FgYellow, color.FgCyan}
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

func (logger *Logger) logOutput(left, right string) {
	if printColor, exists := logger.getColor(left); exists {
		color.Set(printColor, color.Bold)
		fmt.Printf("%s ", logger.pad(left))
		color.Unset()
		color.Set(printColor)
		fmt.Println(right)
		color.Unset()
	} else {
		fmt.Printf("%s %s\n", logger.pad(left), right)
	}
}

// Log logs the given text
func (logger *Logger) Log(role, text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !util.DoesStringArrayContain(logger.SilencedRoles, left) {
			logger.logOutput(left, right)
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
			logger.logOutput(left, right)
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
