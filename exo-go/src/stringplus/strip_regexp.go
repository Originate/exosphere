package stringplus

import "regexp"

// StripRegexp strips off substrings that match the given regex from the
// given text
func StripRegexp(regex, text string) string {
	return string(regexp.MustCompile(regex).ReplaceAll([]byte(text), []byte("")))
}
