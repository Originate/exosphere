package util

import (
	"unicode"
)

// DoesStringArrayContain returns whether the given string slice
// contains the given string.
func DoesStringArrayContain(strings []string, targetString string) bool {
	for _, element := range strings {
		if element == targetString {
			return true
		}
	}
	return false
}

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
//
// From: https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		if unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) {
			out = append(out, unicode.ToLower(runes[i]))
		} else {
			out = append(out, '_')
		}
	}

	return string(out)
}
