package logger

import (
	"fmt"
	"strings"

	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/fatih/color"
	"github.com/willf/pad"
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
	logger := &Logger{Roles: roles, SilencedRoles: silencedRoles, Colors: map[string]color.Attribute{}}
	logger.setColors(roles)
	logger.setLength(roles)
	return logger
}

func (logger *Logger) getColor(role string) (color.Attribute, bool) {
	printColor, exist := logger.Colors[role]
	return printColor, exist
}

func (logger *Logger) getDefaultColors() []color.Attribute {
	return []color.Attribute{color.FgMagenta, color.FgBlue, color.FgYellow, color.FgCyan, color.FgGreen, color.FgWhite}
}

// Log logs the given text
func (logger *Logger) Log(role, text string, trim bool) {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		if !util.DoesStringArrayContain(logger.SilencedRoles, serviceName) {
			logger.logOutput(serviceName, serviceOutput)
		}
	}
}

func (logger *Logger) logOutput(serviceName, serviceOutput string) {
	if printColor, exists := logger.getColor(serviceName); exists {
		regularColor := color.New(printColor)
		boldColor := color.New(printColor, color.Bold)
		fmt.Printf("%s %s\n", boldColor.Sprintf(logger.pad(serviceName)), regularColor.Sprintf(serviceOutput))
	} else {
		boldColor := color.New(color.Bold)
		fmt.Printf("%s %s\n", boldColor.Sprintf(logger.pad(serviceName)), serviceOutput)
	}
}

func (logger *Logger) setColors(roles []string) {
	defaultColors := logger.getDefaultColors()
	for i, role := range roles {
		logger.Colors[role] = defaultColors[i%len(defaultColors)]
	}
}

func (logger *Logger) setLength(roles []string) {
	for _, role := range roles {
		if len(role) > logger.Length {
			logger.Length = len(role)
		}
	}
}

func (logger *Logger) pad(text string) string {
	return pad.Left(text, logger.Length, " ")
}
