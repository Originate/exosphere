package logger

import (
	"fmt"
	"io"
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
	Writer        io.Writer
}

// NewLogger is Logger's constructor
func NewLogger(roles, silencedRoles []string, writer io.Writer) *Logger {
	logger := &Logger{
		Roles:         roles,
		SilencedRoles: silencedRoles,
		Colors:        map[string]color.Attribute{},
		Writer:        writer,
	}
	logger.setColors(roles)
	logger.setLength(roles)
	return logger
}

func (l *Logger) getColor(role string) (color.Attribute, bool) {
	printColor, exist := l.Colors[role]
	return printColor, exist
}

func (l *Logger) getDefaultColors() []color.Attribute {
	return []color.Attribute{color.FgMagenta, color.FgBlue, color.FgYellow, color.FgCyan, color.FgGreen, color.FgWhite}
}

// Log logs the given text
func (l *Logger) Log(role, text string, trim bool) error {
	if trim {
		text = strings.TrimSpace(text)
	}
	for _, line := range strings.Split(text, `\n`) {
		serviceName, serviceOutput := util.ParseDockerComposeLog(role, line)
		if !util.DoesStringArrayContain(l.SilencedRoles, serviceName) {
			err := l.logOutput(serviceName, serviceOutput)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetLogChannel returns a channel which will be endless read from
// and logged with the given role
func (l *Logger) GetLogChannel(role string) chan string {
	textChannel := make(chan string)
	go func() {
		for {
			l.Log(role, <-textChannel, true)
		}
	}()
	return textChannel
}

func (l *Logger) logOutput(serviceName, serviceOutput string) error {
	prefix := color.New(color.Bold).Sprintf(l.pad(serviceName))
	content := serviceOutput
	if printColor, exists := l.getColor(serviceName); exists {
		content = color.New(printColor).Sprintf(serviceOutput)
	}
	_, err := fmt.Fprintf(l.Writer, "%s %s\n", prefix, content)
	return err
}

func (l *Logger) setColors(roles []string) {
	defaultColors := l.getDefaultColors()
	for i, role := range roles {
		l.Colors[role] = defaultColors[i%len(defaultColors)]
	}
}

func (l *Logger) setLength(roles []string) {
	for _, role := range roles {
		if len(role) > l.Length {
			l.Length = len(role)
		}
	}
}

func (l *Logger) pad(text string) string {
	return pad.Left(text, l.Length, " ")
}
