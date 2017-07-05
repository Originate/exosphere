package logger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Originate/exosphere/exo-run-go/helpers"
	"github.com/fatih/color"
)

type Logger struct {
	Roles         []string
	SilencedRoles []string
	Length        int
	Colors        map[string]func(string, ...interface{})
}

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
	logger.SetColors(roles)
	return logger
}

func (logger *Logger) GetColor(role string) (func(string, ...interface{}), bool) {
	chosenColor, exists := logger.Colors[role]
	return chosenColor, exists
}

func (logger *Logger) Log(role, text string, trim bool) {
	// fmt.Println("text", text)
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !helpers.Contains(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", logger.Pad(left), right)
			if color, exists := logger.GetColor(left); exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
}

func (logger *Logger) Error(role, text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		left, right := parseLine(role, line)
		if !helpers.Contains(logger.SilencedRoles, left) {
			output := fmt.Sprintf("%s %s", logger.Pad(left), right)
			if color, exists := logger.GetColor(left); exists {
				color(output)
			} else {
				fmt.Println(output)
			}
		}
	}
}

func (logger *Logger) SetColors(roles []string) {
	defaultColors := []func(string, ...interface{}){color.Magenta, color.Blue, color.Yellow, color.Cyan}
	for i, role := range roles {
		logger.Colors[role] = defaultColors[i%len(defaultColors)]
	}
	for role, _ := range logger.Colors {
		if len(role) > logger.Length {
			logger.Length = len(role)
		}
	}
}

func (logger *Logger) Pad(text string) string {
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
