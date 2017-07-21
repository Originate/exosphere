package processHelpers

import (
	"bufio"
	"io"
)

func readPipe(reader io.Reader, log func(string)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log(scanner.Text())
	}
}
