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

// Logger represents a logger
type Logger struct {
	Length int
	Colors []color.Attribute
	Writer io.Writer
}

// NewLogger is Logger's constructor
func NewLogger(writer io.Writer) *Logger {
	result := &Logger{
		Length: defaultLength,
		Colors: colors,
		Writer: writer,
	}
	return result
}

// Log logs the given text, panics if the write fails
// func (l *Logger) Log(role, text string) {
// 	text = NormalizeDockerComposeLog(text)
// 	for _, line := range strings.Split(text, "\n") {
// 		err := l.logOutput(role, serviceOutput)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// Logf convience method for Log(fmt.Sprintf(...))
func (l *Logger) Logf(role, format string, a ...interface{}) {
	l.Log(role, fmt.Sprintf(format, a...))
}

// Helpers

func (l *Logger) Log(role, serviceOutput string) error {
	printColor := l.getColor(role)
	prefix := color.New(printColor).Sprintf(l.padAndTruncate(role))
	content := strings.TrimSpace(serviceOutput)
	_, err := fmt.Fprintf(l.Writer, "%s %s\n", prefix, content)
	return err
}

func (l *Logger) getColor(role string) color.Attribute {
	hash := fnv.New32()
	hash.Write([]byte(role))
	bitOperand := uint32(1<<31) - 1 // 011111... used to remove signed bit
	index := (hash.Sum32() & bitOperand) % len(l.Colors)
	return l.Colors[index]
}

func (l *Logger) padAndTruncate(text string) string {
	if len(text) > l.Length {
		text = text[:l.Length]
	}
	return pad.Left(text, l.Length, " ")
}
