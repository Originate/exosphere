package execplus

import "bytes"

// Started with the implementation of ScanLines from https://golang.org/src/bufio/scan.go
// and added the part to capture the prompt
func scanLinesOrPrompt(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// If we have a newline, read up to the newline
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// If the chunk ends with ": ", we have a prompt and thus read it all
	if i := bytes.LastIndex(data, []byte(": ")); i >= 0 && i+2 == len(data) {
		return i + 2, dropCR(data[0 : i+2]), nil
	}
	// Request more data.
	return 0, nil, nil
}

// Copied from https://golang.org/src/bufio/scan.go
// needed as a helper for scanLinesOrPrompt
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
