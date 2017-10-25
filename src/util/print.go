package util

import (
	"fmt"
	"strings"
)

// PrintHeader print text within a banner
func PrintHeader(text string) {
	separator := strings.Repeat("*", 80)
	fmt.Printf("\n%s\n%s\n%s\n", separator, text, separator)
}
