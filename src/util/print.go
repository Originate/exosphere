package util

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

// PrintCommandHeader prints a command header
func PrintCommandHeader(writer io.Writer, text string) {
	fmt.Println("")
	_, err := color.New(color.Faint).Fprintln(writer, text)
	if err != nil {
		panic(err)
	}
}

// PrintSectionHeader prints a section header
func PrintSectionHeader(writer io.Writer, text string) {
	fmt.Println("")
	_, err := color.New(color.Underline).Fprint(writer, text)
	if err != nil {
		panic(err)
	}
}

// PrintTimeElapsed prints the time elapsed since startTime
func PrintTimeElapsed(writer io.Writer, command string, startTime time.Time) {
	elapsedTime := time.Since(startTime)
	_, err := color.New(color.Faint).Fprintf(writer, "'%s' elapsed time: %s\n", command, elapsedTime.String())
	if err != nil {
		panic(err)
	}
}

// PrintSectionHeaderf prints a section header with given format
func PrintSectionHeaderf(writer io.Writer, format string, a ...interface{}) {
	PrintSectionHeader(writer, fmt.Sprintf(format, a...))
}
