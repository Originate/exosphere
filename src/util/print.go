package util

import (
	"fmt"
	"io"
	"strings"
)

// PrintHeader prints text within a banner
func PrintHeader(writer io.Writer, text string) {
	separator := strings.Repeat("*", 80)
	fmt.Fprintf(writer, "\n%s\n%s\n%s\n", separator, text, separator)
}

// PrintHeaderf prints formatted text within a banner
func PrintHeaderf(writer io.Writer, format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	separator := strings.Repeat("*", 80)
	fmt.Fprintf(writer, "\n%s\n%s\n%s\n", separator, text, separator)
}
