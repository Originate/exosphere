package stringplus

// StripAnsiColors removes all ANSI color codes from the given string
func StripAnsiColors(text string) string {
	return StripRegexp(`\033\[[0-9;]*m`, text)
}
