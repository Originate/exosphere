package util

import (
	"fmt"
	"hash/fnv"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/willf/pad"
)

var colors = []color.Attribute{color.FgMagenta, color.FgBlue, color.FgYellow, color.FgCyan, color.FgGreen, color.FgWhite}
var defaultLength = 20

// RoleLogger represents a role logger
type RoleLogger struct {
	Length int
	Colors []color.Attribute
	Writer io.Writer
}

// NewRoleLogger is RoleLogger's constructor
func NewRoleLogger(writer io.Writer) *RoleLogger {
	result := &RoleLogger{
		Length: defaultLength,
		Colors: colors,
		Writer: writer,
	}
	return result
}

// Log logs output with a prefixed role, panics if the write fails
func (l *RoleLogger) Log(role, serviceOutput string) {
	printColor, err := l.getColor(role)
	if err != nil {
		panic(err)
	}
	prefix := color.New(printColor).Sprintf(l.padAndTruncate(role))
	content := strings.TrimSpace(serviceOutput)
	_, err = fmt.Fprintf(l.Writer, "%s %s\n", prefix, content)
	if err != nil {
		panic(err)
	}
}

// Logf convience method for Log(fmt.Sprintf(...))
func (l *RoleLogger) Logf(role, format string, a ...interface{}) {
	l.Log(role, fmt.Sprintf(format, a...))
}

// Helpers

func (l *RoleLogger) getColor(role string) (color.Attribute, error) {
	hash := fnv.New32()
	_, err := hash.Write([]byte(role))
	if err != nil {
		return -1, err
	}
	bitOperand := uint32(1<<31) - 1 // 011111... used to remove signed bit
	index := int(hash.Sum32()&bitOperand) % len(l.Colors)
	return l.Colors[index], nil
}

func (l *RoleLogger) padAndTruncate(text string) string {
	if len(text) > l.Length {
		text = text[:l.Length]
	}
	return pad.Left(text, l.Length, " ")
}
