package util

import (
	"fmt"
	"io"
	"strings"
)

// PrintHeader print text within a banner
func PrintHeader(writer io.Writer, text string) {
	separator := strings.Repeat("*", 80)
	fmt.Fprintf(writer, "\n%s\n%s\n%s\n", separator, text, separator)
}
