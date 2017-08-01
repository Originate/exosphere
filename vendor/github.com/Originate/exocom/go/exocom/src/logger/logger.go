package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Originate/exocom/go/exocom/src/client_registry"
	"github.com/Originate/exocom/go/structs"
	"github.com/fatih/color"
)

// Logger prints messages to the a io.Writer
type Logger struct {
	writer io.Writer
}

// NewLogger returns a new Logger
func NewLogger(writer io.Writer) *Logger {
	return &Logger{writer: writer}
}

// Log prints the given text followed by a newline
func (l *Logger) Log(text string) error {
	_, err := fmt.Fprintln(l.writer, text)
	return err
}

// Write prints the given text
func (l *Logger) Write(text string) error {
	_, err := fmt.Fprint(l.writer, text)
	return err
}

// Error prints the given text in red followed by a newline
func (l *Logger) Error(text string) error {
	_, err := color.New(color.FgRed).Fprintln(l.writer, text)
	return err
}

// Warning prints the given text in yellow followed by a newline
func (l *Logger) Warning(text string) error {
	_, err := color.New(color.FgYellow).Fprintln(l.writer, text)
	return err
}

// Header prints the given text in faint (decreased intensity,
// in contrast to bold which is increased intensity) followed by a newline
func (l *Logger) Header(text string) error {
	_, err := color.New(color.Faint).Fprintln(l.writer, text)
	return err
}

// RoutingSetup prints the given routes
func (l *Logger) RoutingSetup(routes clientRegistry.Routes) error {
	err := l.Log("Receiving Routing Setup")
	if err != nil {
		return err
	}
	for service, route := range routes {
		err := l.Log("  --[ " + service + " ]-> " + strings.Join(route.Receives, " + "))
		if err != nil {
			return err
		}
	}
	return nil
}

// Messages prints the given messages and who they were sent to
func (l *Logger) Messages(messages []structs.Message, receivers []string) error {
	for _, message := range messages {
		responseTime := ""
		if message.ResponseTo != "" {
			responseTime = fmt.Sprintf("  ( %v )", message.ResponseTime)
		}
		err := l.Log(message.Sender + "  --[ " + message.Name + " ]->  " + strings.Join(receivers, "and") + responseTime)
		if err != nil {
			return err
		}
		payload, err := json.Marshal(message.Payload)
		if err != nil {
			return err
		}
		err = l.Log(string(payload))
		if err != nil {
			return err
		}
	}
	return nil
}
