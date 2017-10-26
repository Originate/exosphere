package util

import (
	"fmt"
	"io"
	"strings"
)

// PrintBanner prints a banner separator
func PrintBanner(writer io.Writer) {
	separator := strings.Repeat("*", 80)
	fmt.Fprintf(writer, "%s\n", separator)
}
